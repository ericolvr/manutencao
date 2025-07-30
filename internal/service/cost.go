package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type CostService interface {
	List(ctx context.Context) ([]domain.Cost, error)
	FindByID(ctx context.Context, id int) (*domain.Cost, error)
	Update(ctx context.Context, cost *domain.Cost) error
	Delete(ctx context.Context, id int) error
}

type costService struct {
	repo repository.CostRepository
}

func NewCostService(repo repository.CostRepository) CostService {
	return &costService{repo: repo}
}

func (s *costService) List(ctx context.Context) ([]domain.Cost, error) {
	return s.repo.List(ctx)
}

func (s *costService) FindByID(ctx context.Context, id int) (*domain.Cost, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *costService) Update(ctx context.Context, cost *domain.Cost) error {
	return s.repo.Update(ctx, cost)
}

func (s *costService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
