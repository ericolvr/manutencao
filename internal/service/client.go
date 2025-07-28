package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type ClientService interface {
	FindByID(ctx context.Context, id int) (*domain.Client, error)
	Create(ctx context.Context, client *domain.Client) (int, error)
	Update(ctx context.Context, client *domain.Client) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]domain.Client, error)
}

type clientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{repo: repo}
}

func (s *clientService) Create(ctx context.Context, client *domain.Client) (int, error) {
	return s.repo.Create(ctx, client)
}

func (s *clientService) List(ctx context.Context) ([]domain.Client, error) {
	return s.repo.List(ctx)
}

func (s *clientService) FindByID(ctx context.Context, id int) (*domain.Client, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *clientService) Update(ctx context.Context, client *domain.Client) error {
	return s.repo.Update(ctx, client)
}

func (s *clientService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
