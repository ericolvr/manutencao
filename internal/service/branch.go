package service

import (
	"context"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type BranchService interface {
	Create(ctx context.Context, branch *domain.Branch) (int, error)
	List(ctx context.Context) ([]domain.Branch, error)
	FindByID(ctx context.Context, id int) (*domain.Branch, error)
	FindByUniorg(ctx context.Context, uniorg string) (*domain.Branch, error)
	GetByClient(ctx context.Context, client string) ([]domain.Branch, error)
	Update(ctx context.Context, branch *domain.Branch) error
	Delete(ctx context.Context, id int) error
}

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) Create(ctx context.Context, branch *domain.Branch) (int, error) {
	return s.repo.Create(ctx, branch)
}

func (s *branchService) List(ctx context.Context) ([]domain.Branch, error) {
	return s.repo.List(ctx)
}

func (s *branchService) FindByID(ctx context.Context, id int) (*domain.Branch, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *branchService) FindByUniorg(ctx context.Context, uniorg string) (*domain.Branch, error) {
	return s.repo.FindByUniorg(ctx, uniorg)
}

func (s *branchService) GetByClient(ctx context.Context, client string) ([]domain.Branch, error) {
	return s.repo.GetByClient(ctx, client)
}

func (s *branchService) Update(ctx context.Context, branch *domain.Branch) error {
	return s.repo.Update(ctx, branch)
}

func (s *branchService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
