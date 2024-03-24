package service

import (
	"context"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
)

// ProductService handles business logic related to products
type ProductService struct {
	repo repository.Product
}

// NewProductService creates new ProductService with provided repository
func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct creates new product
func (p *ProductService) CreateProduct(ctx context.Context, product *model.Product, categories []model.Category) error {
	return p.repo.CreateProduct(ctx, product, categories)
}

// GetProductsByCategory retrieves products by category name
func (p *ProductService) GetProductsByCategory(ctx context.Context, categoryName string) ([]string, error) {
	return p.repo.GetProductsByCategory(ctx, categoryName)
}

// UpdateProduct updates product and its associated categories
func (p *ProductService) UpdateProduct(ctx context.Context, product *model.Product, categories []model.Category) error {
	return p.repo.UpdateProduct(ctx, product, categories)
}

// DeleteProduct deletes product by its name
func (p *ProductService) DeleteProduct(ctx context.Context, productName string) error {
	productID, err := p.repo.GetProductID(ctx, productName)
	if err != nil {
		return err
	}

	return p.repo.DeleteProduct(ctx, productID)
}
