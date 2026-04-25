package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/golang-api-rest/internal/models"
	"github.com/user/golang-api-rest/internal/services"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{service: s}
}

// GetAll lista todos os usuários.
// @Summary Listar usuários
// @Description Retorna a lista de todos os usuários cadastrados (requer autenticação).
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 401 {string} string "Não autorizado"
// @Router /api/v1/users [get]
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.ListAll()
	if err != nil {
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetByID busca um usuário pelo ID.
// @Summary Buscar usuário por ID
// @Description Retorna os dados de um usuário específico (requer autenticação).
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 200 {object} models.User
// @Failure 400 {string} string "ID inválido"
// @Failure 401 {string} string "Não autorizado"
// @Failure 404 {string} string "Usuário não encontrado"
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Update atualiza os dados de um usuário.
// @Summary Atualizar usuário
// @Description Atualiza nome e/ou email de um usuário (requer autenticação).
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Param request body models.UpdateUserRequest true "Dados para atualizar"
// @Success 200 {object} models.User
// @Failure 400 {string} string "ID ou corpo inválido"
// @Failure 401 {string} string "Não autorizado"
// @Failure 404 {string} string "Usuário não encontrado"
// @Router /api/v1/users/{id} [put]
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	user, err := c.service.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Delete remove um usuário pelo ID.
// @Summary Deletar usuário
// @Description Remove um usuário do sistema (requer role admin).
// @Tags users
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 204 "Sem conteúdo"
// @Failure 400 {string} string "ID inválido"
// @Failure 401 {string} string "Não autorizado"
// @Failure 403 {string} string "Acesso negado"
// @Failure 404 {string} string "Usuário não encontrado"
// @Router /api/v1/users/{id} [delete]
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	if err := c.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
