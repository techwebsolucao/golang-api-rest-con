package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/golang-api-rest/internal/models"
	"github.com/user/golang-api-rest/internal/response"
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
// @Success 200 {object} response.APIResponse "Lista de usuários"
// @Failure 401 {object} response.APIResponse "Não autorizado"
// @Router /users [get]
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.ListAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "erro interno do servidor")
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// GetByID busca um usuário pelo ID.
// @Summary Buscar usuário por ID
// @Description Retorna os dados de um usuário específico (requer autenticação).
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 200 {object} response.APIResponse "Dados do usuário"
// @Failure 400 {object} response.APIResponse "ID inválido"
// @Failure 401 {object} response.APIResponse "Não autorizado"
// @Failure 404 {object} response.APIResponse "Usuário não encontrado"
// @Router /users/{id} [get]
func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id inválido")
		return
	}

	user, err := c.service.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
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
// @Success 200 {object} response.APIResponse "Usuário atualizado com sucesso"
// @Failure 400 {object} response.APIResponse "ID ou corpo inválido"
// @Failure 401 {object} response.APIResponse "Não autorizado"
// @Failure 404 {object} response.APIResponse "Usuário não encontrado"
// @Router /users/{id} [put]
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id inválido")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	user, err := c.service.Update(id, &req)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// Delete remove um usuário pelo ID.
// @Summary Deletar usuário
// @Description Remove um usuário do sistema (requer role admin).
// @Tags users
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 200 {object} response.APIResponse "Usuário deletado com sucesso"
// @Failure 400 {object} response.APIResponse "ID inválido"
// @Failure 401 {object} response.APIResponse "Não autorizado"
// @Failure 403 {object} response.APIResponse "Acesso negado"
// @Failure 404 {object} response.APIResponse "Usuário não encontrado"
// @Router /users/{id} [delete]
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id inválido")
		return
	}

	if err := c.service.Delete(id); err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Message(w, http.StatusOK, "usuário deletado com sucesso")
}
