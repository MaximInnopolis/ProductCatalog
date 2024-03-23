package repository

import (
	"context"
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
)

type Authorization interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, user *model.User) (*model.User, error)
}

type Category interface {
	CreateCategory(ctx context.Context, category *model.Category) (int64, error)
	GetAllCategories(ctx context.Context) ([]string, error)
	UpdateCategory(ctx context.Context, categoryName string, category *model.Category) error
	DeleteCategory(ctx context.Context, categoryName string) error
}

type Repository struct {
	Authorization
	Category
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthDatabase(db),
		Category:      NewCategoryDatabase(db),
	}
}
