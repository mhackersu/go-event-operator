package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

type InputData struct {
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
}

type CalculationResult struct {
	Symbol string  `json:"symbol"`
	Result float64 `json:"result"`
}

func fetchDataFromGCP(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	reader, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data := make([]byte, reader.Attrs.Size)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func performCalculations(data []InputData) []CalculationResult {
	results := []CalculationResult{}
	for _, entry := range data {
		results = append(results, CalculationResult{
			Symbol: entry.Symbol,
			Result: entry.Amount * 1.05, // Example calculation (e.g., 5% increase)
		})
	}
	return results
}

func main() {
	ctx := context.Background()
	bucket := os.Getenv("GCP_BUCKET")
	object := os.Getenv("GCP_OBJECT")

	data, err := fetchDataFromGCP(ctx, bucket, object)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	var inputs []InputData
	if err := json.Unmarshal(data, &inputs); err != nil {
		log.Fatalf("Failed to parse input data: %v", err)
	}

	results := performCalculations(inputs)
	for _, result := range results {
		jsonData, _ := json.Marshal(result)
		fmt.Println(string(jsonData))
	}
}
