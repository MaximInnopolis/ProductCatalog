# ProductCatalog
API для каталога товаров

get request:
curl http://localhost:8080/categories

post request:
curl -X POST \
-H "Content-Type: application/json" \
-d '{"username": "exampleuser", "password": "examplepassword"}' \
http://localhost:8080/register

/
├── .github
│   └── workflows
│       └── main.yaml
├── .gitignore
├── /cmd
│   └── /catalog_api_server
│       └── main.go
├── /internal
│   ├── /api
│   │   ├── api.go
│   │   └── handlers.go
│   ├── /auth
│   │   └── auth.go
│   ├── /database
│   │   └── database.go
│   ├── /logger
│   │   └── logger.go
│   └── /models
│       ├── category.go
│       ├── product.go
│       └── user.go
├── /scripts
│   └── /database
│       └── migrate.go
├── /tests
│   ├── database_test.go
│   └── handlers_test.go
├── /docs
│   └── README.md
├── /data
│   └── database.db
├── go.mod
└── go.sum

