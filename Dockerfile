# Этап сборки
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/api/main.go

# Финальный легкий образ
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный файл из этапа сборки
COPY --from=builder /app/bot .

# Запускаем приложение
CMD ["./bot"]