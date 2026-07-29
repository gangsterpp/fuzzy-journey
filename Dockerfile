FROM golang:1.26-alpine AS builder

WORKDIR /app

# Разрешаем Go автоматически подтягивать нужную версию toolchain из go.mod, если она новее
ENV GOTOOLCHAIN=auto

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

EXPOSE 8080

# Запускаем приложение
CMD ["./bot"]