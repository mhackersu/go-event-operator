package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type InputData struct {
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
}

type CalculationResult struct {
	Symbol string  `json:"symbol"`
	Result float64 `json:"result"`
}

func fetchDataFromPubSub(ctx context.Context, projectID, subscriptionID string) ([]InputData, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	sub := client.Subscription(subscriptionID)
	var inputData []InputData

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		var data InputData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Failed to parse message: %v", err)
			msg.Nack()
			return
		}
		inputData = append(inputData, data)
		msg.Ack()
	})

	if err != nil {
		return nil, err
	}

	return inputData, nil
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

func publishResultsToPubSub(ctx context.Context, projectID, topicID string, results []CalculationResult) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	topic := client.Topic(topicID)
	for _, result := range results {
		jsonData, _ := json.Marshal(result)
		res := topic.Publish(ctx, &pubsub.Message{Data: jsonData})
		_, err := res.Get(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT")
	subscriptionID := os.Getenv("GCP_SUBSCRIPTION")
	topicID := os.Getenv("GCP_TOPIC")

	inputs, err := fetchDataFromPubSub(ctx, projectID, subscriptionID)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	results := performCalculations(inputs)
	if err := publishResultsToPubSub(ctx, projectID, topicID, results); err != nil {
		log.Fatalf("Failed to publish results: %v", err)
	}

	if os.Getenv("PRINT_JSONL") == "true" {
		for _, result := range results {
			jsonData, _ := json.Marshal(result)
			fmt.Println(string(jsonData))
		}
	}
}
