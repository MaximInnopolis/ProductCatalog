![workflow](https://github.com/MaximInnopolis/ProductCatalog/actions/workflows/main.yaml/badge.svg)
# Инструкция по запуску проекта

Этот файл содержит инструкции по запуску проекта с использованием Dockerfile. Проект представляет собой API для каталога товаров. В проекте используется SQLite.

## Предварительные требования

Убедитесь, что у вас установлен Docker и Docker Compose.

## Шаги запуска проекта

1. **Получение образа из Docker Hub:**

   Выполните команду для загрузки образа проекта из Docker Hub:

   ```bash
   docker pull madfisher/catalog-api:latest
    ```
2. **Запуск контейнера:**

   Выполните команду для запуска контейнера из загруженного образа:

   ```bash
    docker run -d -p 8080:8080 madfisher/catalog-api:latest
    ```
   Где:
   1. `-d`: запуск контейнера в фоновом режиме.
   2. `-p 8080:8080`: проброс портов, где первый номер порта (8080) - порт вашего хоста, а второй номер порта (8080) - порт внутри контейнера, на котором работает приложение.

    Теперь API доступно по адресу http://localhost:8080.

## Время, затраченное на разработку каждой части проекта

Суммарно 10 часов:

- Авторизация: 2 часа
- Работа с базой данных (добавление, обновление, удаление, вывод записей): 2 часа
- Установление иерархии проекта, проектирование логики CI, описание моделей: 40 минут
- Настройка DockerFile, добавление Docker образа на Dockerhub: 20 минут
- Написание тестов: 1,5 часа
- Разработка прочих хэндлеров: 1 час
- Написание logger, миграции и создание тестовых записей в таблице: 40 минут
- Сборка товаров из внешнего источника: 30 минут
- Написание README.md и комментарии к коду - 20 минут

## Описание API

API предоставляет следующие методы:

- **POST /auth/register:** Зарегистрировать нового пользователя.
- **POST /auth/login:** Авторизовать пользователя.


- **GET /categories/list:** Получить список категорий.
- **POST /categories/new:** Создать новую категорию (требуется авторизация).
- **PUT /categories/{name}:** Обновить существующую категорию (требуется авторизация).
- **DELETE /categories/{name}:** Удалить категорию по имени (требуется авторизация).


- **GET /products/{name}:** Получить список товаров в указанной категории.
- **POST /products/new:** Добавить новый товар (требуется авторизация).
- **PUT /products:** Обновить существующий товар (требуется авторизация).
- **DELETE /products/{name}:** Удалить товар по имени (требуется авторизация).

### Примеры запросов

Примеры curl-запросов для каждого из вышеперечисленных методов API:

#### Зарегистрировать нового пользователя
```bash
curl -X POST \
-H "Content-Type: application/json" \
-d '{"username": "exampleuser", "password": "examplepassword"}' \
http://localhost:8080/auth/register
```
#### Авторизовать пользователя
```bash
curl -X POST \
-H "Content-Type: application/json" \
-d '{"username": "exampleuser", "password": "examplepassword"}' \
http://localhost:8080/auth/login
```


#### Получить список категорий
```bash
curl http://localhost:8080/categories/list
```

#### Создать новую категорию товаров
```bash
curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_TOKEN" \
http://localhost:8080/categories/new
```

#### Обновить существующую категорию товаров
```bash
curl -X PUT -H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_TOKEN" \
-d '{"Name": "New Category Name"}' \
http://localhost:8080/categories/CategoryName
```

#### Удалить категорию товаров по имени
```bash
curl -X DELETE \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_TOKEN" \
http://localhost:8080/categories/CategoryName
```

#### Получить список товаров в указанной категории
```bash
curl http://localhost:8080/products/CategoryName
```

#### Добавить новый товар в указанную категорию
```bash
curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_TOKEN" \
-d '{"Name": "New Product Name"}' \
http://localhost:8080/products/new
```

#### Обновить товар
```bash
curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_TOKEN" \
-d '{"Name": "Product Name", "Category": "Category Name"}' \
http://localhost:8080/products
```

#### Удалить товар по имени
```bash
curl -X DELETE \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_TOKEN" \
http://localhost:8080/products/ProductName
```

## Структура проекта
```
/ProductCatalog
│
├── /.github
│   └── workflows
│       └── main.yaml
│
├── .gitignore
├── /cmd
│   └── /catalog_api_server
│       └── main.go
│ 
├── /internal
│   ├── /handler
│   │   ├── auth.go
│   │   ├── category_handlers.go
│   │   ├── product_handlers.go
│   │   ├── middleware.go
│   │   └── handlers.go
│   │
│   ├── /repository
│   │   ├── auth_database.go
│   │   ├── category_database.go
│   │   ├── product_database.go
│   │   ├── repository.go
│   │   └── database.go
│   │
│   ├── /service
│   │   ├── auth.go
│   │   ├── category.go
│   │   ├── product.go
│   │   └── service.go
│   │
│   ├── /logger
│   │   └── logger.go
│   │
│   ├── /utils
│   │   ├── response_writer.go
│   │   └── env_loader.go
│   │
│   └── /model
│       ├── category.go
│       ├── product.go
│       ├── response.go
│       └── user.go
│
├── /scripts
│   ├── migrate.go
│   ├── create_records.go
│   └── data_collection.go
│
├── /tests
│   ├── /api_test
│   │   └── handlers_test.go
│   ├── /auth_test
│   │   └── auth_test.go
│   ├── /database
│   │   └── database_test.go
│   ├── /models_test
│   │   ├── category_test.go
│   │   ├── product_test.go
│   │   └── user_test.go
│   └── /scripts_test
│       ├── migrate_test.go
│       ├── data_collection.go
│       └── create_records.go
│
├── /docs
│   ├── Feedback.md
│   ├── README.md 
│   └── Task.md
│
├── /data
│   └── database.db
│
├── go.mod
├── go.sum
├── .env
├── .dockerignore
└── Dockerfile
```


