FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
# Копируем .env внутрь контейнера чтобы viper мог его прочитать
COPY .env .env

EXPOSE 8080

CMD ["./server"]
