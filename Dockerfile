FROM golang:1.21 AS builder

# Устанавливаем переменную окружения для отключения кэша Go
ENV GO111MODULE=on

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копирование зависимостей проекта
COPY go.mod go.sum ./

# Загрузка зависимости
RUN go mod download

# Копирование файлов проекта внутрь контейнера
COPY . .

# Сборка приложения внутри контейнера
RUN go build -o catalog_api_server ./cmd/catalog_api_server/main.go

# Запуск приложения при запуске контейнера
CMD ["./catalog_api_server"]
