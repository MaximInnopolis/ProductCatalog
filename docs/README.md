# ProductCatalog
API для каталога товаров

get categories request:
curl http://localhost:8080/categories/list

post category request:
curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTAwMTkyNzQsImlkIjoxLCJzdWIiOiJleGFtcGxldXNlciJ9.vWdoR9zqir_dYcH25Bqp8U7GTDujaeJbXup2zR4InS8" \
http://localhost:8080/categories/new

update category request:
curl -X PUT \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTAwMTgzNDYsImlkIjoxLCJzdWIiOiJleGFtcGxldXNlciJ9.olmHkaiko8urYeN5jpsnL1ChldJcdXkZSOd_IC3RTDQ" \
-d '{ "Name": "Fish"}' \
http://localhost:8080/categories/new

delete categories request:
curl -X DELETE \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTAwMTgzNDYsImlkIjoxLCJzdWIiOiJleGFtcGxldXNlciJ9.olmHkaiko8urYeN5jpsnL1ChldJcdXkZSOd_IC3RTDQ" \
http://localhost:8080/categories/Food

get products by category request:
curl http://localhost:8080/products/smileys%20and%20people

delete products request:
curl -X DELETE \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTAwMjIxNjYsImlkIjoxLCJzdWIiOiJleGFtcGxldXNlciJ9.330o7IyhQFWl-9TNHJslum1A_R6tW_q0e8AqnaV6n2w" \
http://localhost:8080/products/Rick


post register request:
curl -X POST \
-H "Content-Type: application/json" \
-d '{"username": "exampleuser", "password": "examplepassword"}' \
http://localhost:8080/auth/register


post login request:
curl -X POST \
-H "Content-Type: application/json" \
-d '{"username": "exampleuser", "password": "examplepassword"}' \
http://localhost:8080/auth/login

go test -v ./tests/scripts


/ProductCatalog
    │
    ├── .github
    │   └── workflows
    │       └── main.yaml
    ├── .gitignore
    ├── /cmd
    │   └── /catalog_api_server
    │       └── main.go
    ├── /internal
    │   ├── /api
    │   │   ├── category_handlers.go
    │   │   ├── product_handlers.go
    │   │   └── handlers.go
    │   ├── /auth
    │   │   └── auth.go
    │   ├── /database
    │   │   └── database.go
    │   ├── /logger
    │   │   └── logger.go
    │   ├── /utils
    │   │   └── response_writer.go
    │   │   └── env_loader.go
    │   └── /models
    │       ├── category.go
    │       ├── product.go
    │       ├── product_database.go
    │       ├── product_helpers.go
    │       ├── response.go
    │       └── user.go
    ├── /scripts
    │   ├── migrate.go
    │   └── create_records.go
    │   └── data_collection.go
    ├── /tests
    │   ├── /api
    │   │   └── handlers_test.go
    │   ├── /auth
    │   │   └── auth_test.go
    │   ├── /database
    │   │   └── database_test.go
    │   ├── /models
    │   │   ├── category_test.go
    │   │   └── product_test.go
    │   │   └── user_test.go
    │   └── /scripts
    │       └── migrate_test.go
    ├── /docs
    │   └── README.md
    ├── /data
    │   └── database.db
    ├── go.mod
    └── go.sum
    └── .env
    └── .dockerignore
    └── Dockerfile

