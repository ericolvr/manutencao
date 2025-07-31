package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type DistanceService interface {
	Create(ctx context.Context, distance *domain.Distance) (int, error)
	List(ctx context.Context) ([]domain.Distance, error)
	FindByID(ctx context.Context, id int) (*domain.Distance, error)
	FindByNumber(ctx context.Context, number string) (*domain.Distance, error)
	Update(ctx context.Context, distance *domain.Distance) error
	Delete(ctx context.Context, id int) error
}

type distanceService struct {
	repo repository.DistanceRepository
}

func NewDistanceService(repo repository.DistanceRepository) DistanceService {
	return &distanceService{repo: repo}
}

func (s *distanceService) Create(ctx context.Context, distance *domain.Distance) (int, error) {
	return s.repo.Create(ctx, distance)
}

func (s *distanceService) List(ctx context.Context) ([]domain.Distance, error) {
	return s.repo.List(ctx)
}

func (s *distanceService) FindByID(ctx context.Context, id int) (*domain.Distance, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *distanceService) FindByNumber(ctx context.Context, number string) (*domain.Distance, error) {
	return s.repo.FindByNumber(ctx, number)
}

func (s *distanceService) Update(ctx context.Context, distance *domain.Distance) error {
	return s.repo.Update(ctx, distance)
}

func (s *distanceService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
