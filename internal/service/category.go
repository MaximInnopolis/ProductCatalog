package service

import (
	"context"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
)

// CategoryService handles business logic related to categories
type CategoryService struct {
	repo repository.Category
}

// NewCategoryService creates new CategoryService with the provided repository
func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategory creates new category
func (s *CategoryService) CreateCategory(ctx context.Context, category *model.Category) (int64, error) {
	return s.repo.CreateCategory(ctx, category)
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetAllCategories(ctx)
}

// UpdateCategory updates category by its name
func (s *CategoryService) UpdateCategory(ctx context.Context, categoryName string, category *model.Category) error {
	return s.repo.UpdateCategory(ctx, categoryName, category)
}

// DeleteCategory deletes category by its name
func (s *CategoryService) DeleteCategory(ctx context.Context, categoryName string) error {
	return s.repo.DeleteCategory(ctx, categoryName)
}
