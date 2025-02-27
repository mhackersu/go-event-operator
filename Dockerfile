FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o golden_service

FROM gcr.io/distroless/base-debian10
WORKDIR /root/
COPY --from=builder /app/golden_service ./
ENTRYPOINT ["./golden_service"]
