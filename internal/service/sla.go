package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type SlaService interface {
	Create(ctx context.Context, sla *domain.Sla) (int, error)
	List(ctx context.Context) ([]domain.Sla, error)
	FindByID(ctx context.Context, id int) (*domain.Sla, error)
	GetByClient(ctx context.Context, client_id int) ([]domain.Sla, error)
	Update(ctx context.Context, sla *domain.Sla) error
	Delete(ctx context.Context, id int) error
}

type slaService struct {
	repo repository.SlaRepository
}

func NewSlaService(repo repository.SlaRepository) SlaService {
	return &slaService{repo: repo}
}

func (s *slaService) Create(ctx context.Context, sla *domain.Sla) (int, error) {
	return s.repo.Create(ctx, sla)
}

func (s *slaService) List(ctx context.Context) ([]domain.Sla, error) {
	return s.repo.List(ctx)
}

func (s *slaService) FindByID(ctx context.Context, id int) (*domain.Sla, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *slaService) GetByClient(ctx context.Context, client_id int) ([]domain.Sla, error) {
	return s.repo.GetByClient(ctx, client_id)
}

func (s *slaService) Update(ctx context.Context, sla *domain.Sla) error {
	return s.repo.Update(ctx, sla)
}

func (s *slaService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
