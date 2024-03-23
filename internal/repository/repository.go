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

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthDatabase(db),
	}
}
