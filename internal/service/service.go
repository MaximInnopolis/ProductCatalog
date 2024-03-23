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

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
