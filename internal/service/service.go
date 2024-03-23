package service

import (
	"context"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user *model.User) error
	GenerateToken(ctx context.Context, user *model.User) (string, error)
	IsTokenValid(ctx context.Context, tokenString string) (bool, error)
}

type Category interface {
	CreateCategory(ctx context.Context, category *model.Category) (int64, error)
	GetAllCategories(ctx context.Context) ([]string, error)
	UpdateCategory(ctx context.Context, categoryName string, category *model.Category) error
	DeleteCategory(ctx context.Context, categoryName string) error
}

type Product interface {
	CreateProduct(ctx context.Context, product *model.Product, categories []model.Category) error
	GetProductsByCategory(ctx context.Context, categoryName string) ([]string, error)
	UpdateProduct(ctx context.Context, product *model.Product, categories []model.Category) error
	DeleteProduct(ctx context.Context, productName string) error
}

type Service struct {
	Authorization
	Category
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Category:      NewCategoryService(repos.Category),
		Product:       NewProductService(repos.Product),
	}
}