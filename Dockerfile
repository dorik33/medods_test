FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main
# Устанавливаем goose для выполнения миграций
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY .env .

EXPOSE 8080
# Выполняем миграции перед запуском приложения
CMD ["sh", "-c", "goose -dir ./migrations postgres \"$DATABASE_URL\" up && ./main"]