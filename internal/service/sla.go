package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type SlaService interface {
	Create(ctx context.Context, sla *domain.Sla) (int, error)
	List(ctx context.Context) ([]domain.Sla, error)
	ListWithClientNames(ctx context.Context) ([]dto.SlaResponse, error)
	FindByID(ctx context.Context, id int) (*domain.Sla, error)
	FindByParams(ctx context.Context, client_id int, priority int) (*dto.SlaJustHours, error)
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

func (s *slaService) ListWithClientNames(ctx context.Context) ([]dto.SlaResponse, error) {
	slaWithClients, err := s.repo.ListWithClientNames(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.SlaResponse
	for _, sla := range slaWithClients {
		response := dto.SlaResponse{
			ID:         sla.ID,
			ClientID:   sla.ClientID,
			ClientName: sla.ClientName,
			Priority:   sla.Priority,
			Hours:      sla.Hours,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *slaService) FindByID(ctx context.Context, id int) (*domain.Sla, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *slaService) FindByParams(ctx context.Context, client_id int, priority int) (*dto.SlaJustHours, error) {
	return s.repo.FindByParams(ctx, client_id, priority)
}
func (s *slaService) Update(ctx context.Context, sla *domain.Sla) error {
	return s.repo.Update(ctx, sla)
}

func (s *slaService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
