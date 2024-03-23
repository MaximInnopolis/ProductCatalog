package service

import (
	"context"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) CreateProduct(ctx context.Context, product *model.Product, categories []model.Category) error {
	return p.repo.CreateProduct(ctx, product, categories)
}

func (p *ProductService) GetProductsByCategory(ctx context.Context, categoryName string) ([]string, error) {
	return p.repo.GetProductsByCategory(ctx, categoryName)
}

func (p *ProductService) UpdateProduct(ctx context.Context, product *model.Product, categories []model.Category) error {
	return p.repo.UpdateProduct(ctx, product, categories)
}

func (p *ProductService) DeleteProduct(ctx context.Context, productName string) error {
	return p.repo.DeleteProduct(ctx, productName)
}
