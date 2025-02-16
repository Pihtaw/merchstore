# 1. Используем официальный образ Go для сборки
FROM golang:1.21 AS builder

# 2. Устанавливаем рабочую директорию в /merch_store
WORKDIR /merch_store

# 3. Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# 4. Копируем весь код
COPY . .

# 5. Сборка приложения
RUN go build -o merch_store

# 6. Используем минимальный образ Debian
FROM debian:bookworm-slim

# 7. Устанавливаем рабочую директорию
WORKDIR /app

# 8. Копируем бинарник из builder-стадии с правильным путем
COPY --from=builder /merch_store/merch_store .

# 9. Открываем порт (если приложение использует HTTP)
EXPOSE 8080

# 10. Запускаем приложение
CMD ["./merch_store"]
