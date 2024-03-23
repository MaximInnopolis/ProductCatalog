package service

import (
	"context"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
)

type CategoryService struct {
	repo repository.Category
}

func (s CategoryService) CreateCategory(ctx context.Context, category *model.Category) (int64, error) {
	return s.repo.CreateCategory(ctx, category)
}

func (s CategoryService) GetAllCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s CategoryService) UpdateCategory(ctx context.Context, categoryName string, category *model.Category) error {
	return s.repo.UpdateCategory(ctx, categoryName, category)
}

func (s CategoryService) DeleteCategory(ctx context.Context, categoryName string) error {
	return s.repo.DeleteCategory(ctx, categoryName)
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}
