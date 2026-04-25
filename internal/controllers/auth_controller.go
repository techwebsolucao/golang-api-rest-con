package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/user/golang-api-rest/internal/models"
	"github.com/user/golang-api-rest/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register cria uma nova conta de usuário.
// @Summary Registrar novo usuário
// @Description Cria um novo usuário e retorna os dados + tokens JWT. Um email de verificação é enviado.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Dados do registro"
// @Success 201 {object} map[string]interface{} "user + token"
// @Failure 400 {string} string "Erro de validação"
// @Failure 409 {string} string "Email já cadastrado"
// @Router /api/v1/auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	user, token, err := c.authService.Register(&req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email já cadastrado" {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// Login autentica o usuário e retorna tokens JWT.
// @Summary Login
// @Description Autentica com email e senha e retorna access + refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciais"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {string} string "Credenciais inválidas"
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// Refresh renova o access token usando um refresh token válido.
// @Summary Renovar tokens
// @Description Gera novos access + refresh tokens a partir de um refresh token válido.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {string} string "Refresh token inválido ou expirado"
// @Router /api/v1/auth/refresh [post]
func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	token, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// Verify confirma o email do usuário através de um token de verificação.
// @Summary Verificar email
// @Description Marca o email do usuário como verificado usando o token enviado por email.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string false "Token de verificação"
// @Param request body models.VerifyEmailRequest false "Token no body (alternativa)"
// @Success 200 {object} map[string]string "mensagem de sucesso"
// @Failure 400 {string} string "Token inválido ou expirado"
// @Router /api/v1/auth/verify-email [get]
// @Router /api/v1/auth/verify-email [post]
func (c *AuthController) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		var req models.VerifyEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
			http.Error(w, "token é obrigatório (query param ?token= ou body JSON)", http.StatusBadRequest)
			return
		}
		token = req.Token
	}

	if err := c.authService.VerifyEmail(token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "email verificado com sucesso"})
}
