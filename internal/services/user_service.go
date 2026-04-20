package services

import (
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
