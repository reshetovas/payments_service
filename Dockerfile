FROM golang:1.23-alpine

# Установка зависимостей
RUN apk update && apk add --no-cache \
    git \
    sqlite \
    gcc \
    musl-dev \
    make \
    bash \
    curl

# Установка goose
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

# Создание рабочей директории
WORKDIR /app

# Копируем go-модули и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Сборка проекта
RUN go build -o payments_service .

# Запуск приложения
CMD ["./payments_service"]