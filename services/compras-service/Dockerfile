FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o procurement-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/procurement-service .
EXPOSE 8085
CMD ["./procurement-service"]

