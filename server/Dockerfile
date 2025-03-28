# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 9090
CMD ["./main"]
