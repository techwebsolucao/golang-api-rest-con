package services

import (
	"fmt"

	"github.com/user/golang-api-rest/internal/models"
	"github.com/user/golang-api-rest/internal/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListAll() ([]models.User, error) {
	return s.repo.All()
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado: %w", err)
	}
	return user, nil
}

func (s *UserService) Update(id int, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado: %w", err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	return user, nil
}

func (s *UserService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}