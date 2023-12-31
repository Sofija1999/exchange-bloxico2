package usecases

import (
	"context"

	domain "github.com/Bloxico/exchange-gateway/sofija/core/domain"
	ports "github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/pkg/errors"
)

// Check this service satisfies interface
var _ ports.EgwProductUsecase = (*EgwProductService)(nil)

type EgwProductService struct {
	productRepo *repo.EgwProductRepository
}

func NewEgwProductService(productRepo *repo.EgwProductRepository) *EgwProductService {
	return &EgwProductService{
		productRepo: productRepo,
	}
}

func (s *EgwProductService) InsertProduct(ctx context.Context, product *domain.EgwProduct) error {
	err := s.productRepo.Insert(ctx, product)
	if err != nil {
		return errors.Wrap(err, "Failed to insert product")
	}
	return nil
}

func (s *EgwProductService) FindByID(ctx context.Context, id string) (*domain.EgwProduct, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve product")
	}

	return product, nil
}

func (s *EgwProductService) Update(ctx context.Context, egwProduct *domain.EgwProduct) error {
	err := s.productRepo.Update(ctx, egwProduct)
	if err != nil {
		return err
	}
	return nil
}

func (s *EgwProductService) Delete(ctx context.Context, id string) error {
	err := s.productRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
