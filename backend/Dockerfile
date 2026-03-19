FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM alpine:latest
WORKDIR /app
# Копируем ТОЛЬКО бинарник
COPY --from=builder /app/server .

# .env НЕ КОПИРУЕМ. Данные придут через Environment Variables
EXPOSE 8080
CMD ["./server"]