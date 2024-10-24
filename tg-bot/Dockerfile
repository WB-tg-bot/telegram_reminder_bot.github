# Используем официальный образ Go как базовый
FROM golang:latest

# Устанавливаем временную зону
RUN apt-get update && apt-get install -y tzdata

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Строим бинарный файл
RUN go build -o tg-bot ./cmd/main.go

# Указываем команду для запуска контейнера
CMD ["./tg-bot"]