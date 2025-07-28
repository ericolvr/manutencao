package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type ProviderService interface {
	Create(ctx context.Context, provider *domain.Provider) (int, error)
	List(ctx context.Context) ([]domain.Provider, error)
	FindByID(ctx context.Context, id int) (*domain.Provider, error)
	FindByName(ctx context.Context, name string) (*domain.Provider, error)
	Update(ctx context.Context, provider *domain.Provider) error
	Delete(ctx context.Context, id int) error
}

type providerService struct {
	repo repository.ProviderRepository
}

func NewProviderService(repo repository.ProviderRepository) ProviderService {
	return &providerService{repo: repo}
}

func (s *providerService) Create(ctx context.Context, provider *domain.Provider) (int, error) {
	return s.repo.Create(ctx, provider)
}

func (s *providerService) List(ctx context.Context) ([]domain.Provider, error) {
	return s.repo.List(ctx)
}

func (s *providerService) FindByID(ctx context.Context, id int) (*domain.Provider, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *providerService) FindByName(ctx context.Context, name string) (*domain.Provider, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *providerService) Update(ctx context.Context, provider *domain.Provider) error {
	return s.repo.Update(ctx, provider)
}

func (s *providerService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
