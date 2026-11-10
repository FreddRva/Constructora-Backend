FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/inventory-service .
EXPOSE 8082
CMD ["./inventory-service"]

