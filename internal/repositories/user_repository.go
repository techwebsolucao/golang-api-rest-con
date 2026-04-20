package repositories

import (
	"errors"
	"github.com/user/golang-api-rest/internal/models"
)

type UserRepository interface {
	GetByID(id int) (*models.User, error)
	All() ([]models.User, error)
}

type MemoryRepository struct {
	users []models.User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: []models.User{
			{ID: 1, Name: "João Silva", Email: "joao@example.com"},
			{ID: 2, Name: "Maria Souza", Email: "maria@example.com"},
		},
	}
}

func (r *MemoryRepository) All() ([]models.User, error) {
	return r.users, nil
}

func (r *MemoryRepository) GetByID(id int) (*models.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}
