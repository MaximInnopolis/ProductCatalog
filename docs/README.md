# ProductCatalog
API для каталога товаров

get categories request:
curl http://localhost:8080/categories/list

post category request:
curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6IjIifQ.ABzj_AVpr2hTZm4EWtuNbQ7kMIk8Gel32z7fMuLth24" \
-d '{ "Name": "Fish"}' \
http://localhost:8080/categories/new


update category request:
curl -X PUT \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6IjIifQ.ABzj_AVpr2hTZm4EWtuNbQ7kMIk8Gel32z7fMuLth24" \
-d '{ "Name": "Fish"}' \
http://localhost:8080/categories/Fish


delete request:

curl -X DELETE \ 
-H "Content-Type: application/json" \ 
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6IjIifQ.ABzj_AVpr2hTZm4EWtuNbQ7kMIk8Gel32z7fMuLth24" \
http://localhost:8080/products/Rick



post register request:
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
│   │   ├── category_handlers.go
│   │   ├── product_handlers.go
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
│       ├── product_database.go
│       ├── product_helpers.go
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

