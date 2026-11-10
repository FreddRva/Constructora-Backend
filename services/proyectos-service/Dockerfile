FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o project-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/project-service .
EXPOSE 8081
CMD ["./project-service"]

