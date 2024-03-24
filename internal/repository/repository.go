package repository

import (
	"context"
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
)

// Authorization interface defines methods for user authorization
type Authorization interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, user *model.User) (*model.User, error)
}

// Category interface defines methods for category management
type Category interface {
	CreateCategory(ctx context.Context, category *model.Category) (int64, error)
	GetAllCategories(ctx context.Context) ([]string, error)
	UpdateCategory(ctx context.Context, categoryName string, category *model.Category) error
	DeleteCategory(ctx context.Context, categoryName string) error
}

// Product interface defines methods for product management
type Product interface {
	CreateProduct(ctx context.Context, product *model.Product, categories []model.Category) error
	GetProductsByCategory(ctx context.Context, categoryName string) ([]string, error)
	UpdateProduct(ctx context.Context, product *model.Product, categories []model.Category) error
	DeleteProduct(ctx context.Context, productId int64) error
	GetCategoryID(ctx context.Context, categoryName string) (int64, error)
	GetProductID(ctx context.Context, productName string) (int64, error)
}

type Repository struct {
	Authorization
	Category
	Product
}

// NewRepository creates new repository with given database connection
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthDatabase(db),
		Category:      NewCategoryDatabase(db),
		Product:       NewProductDatabase(db),
	}
}
