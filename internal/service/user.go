package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) (int, error)
	List(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	FindByName(ctx context.Context, name string) ([]domain.User, error)
	FindByMobile(ctx context.Context, mobile string) ([]domain.User, error)
	FindUsersToTicket(ctx context.Context) ([]dto.UsersToTicket, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	Authenticate(ctx context.Context, mobile, password string) (*dto.AuthResponse, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewUserService(repo repository.UserRepository, jwtSecret []byte) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (s *userService) Create(ctx context.Context, user *domain.User) (int, error) {
	hashedPassword, err := domain.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	user.Password = hashedPassword

	return s.repo.Create(ctx, user)
}

func (s *userService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) FindByID(ctx context.Context, id int) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) FindByName(ctx context.Context, name string) ([]domain.User, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *userService) FindByMobile(ctx context.Context, mobile string) ([]domain.User, error) {
	return s.repo.FindByMobile(ctx, mobile)
}

func (s *userService) FindUsersToTicket(ctx context.Context) ([]dto.UsersToTicket, error) {
	return s.repo.FindUsersToTicket(ctx)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	if user.Password != "" {
		hashedPassword, err := domain.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("error generating hash: %w", err)
		}
		user.Password = hashedPassword
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Authenticate(ctx context.Context, mobile, password string) (*dto.AuthResponse, error) {
	users, err := s.repo.FindByMobile(ctx, mobile)
	if err != nil {
		fmt.Printf("Erro ao buscar usuário: %v\n", err)
		return nil, errors.New("invalid mobile or password")
	}

	if len(users) == 0 {
		fmt.Printf("Nenhum usuário encontrado com o mobile: %s\n", mobile)
		return nil, errors.New("invalid mobile or password")
	}

	user := users[0]

	if !user.Status {
		return nil, errors.New("user account is inactive")
	}

	fmt.Printf("Usuário encontrado: ID=%d, Nome=%s, Mobile=%s, Role=%d\n",
		user.ID, user.Name, user.Mobile, user.Role)

	err = domain.CheckPassword(user.Password, password)
	if err != nil {
		return nil, errors.New("invalid mobile or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":   user.Name,
		"mobile": user.Mobile,
		"role":   user.Role,
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &dto.AuthResponse{
		Name:  user.Name,
		Token: tokenString,
		Role:  user.Role,
	}, nil
}
