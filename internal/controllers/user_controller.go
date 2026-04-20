package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/user/golang-api-rest/internal/services"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.ListAll()
	
	if err != nil {
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
