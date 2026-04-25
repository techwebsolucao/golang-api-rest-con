package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/user/golang-api-rest/internal/models"
	"github.com/user/golang-api-rest/internal/response"
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
// @Success 201 {object} response.APIResponse "Usuário criado com sucesso"
// @Failure 400 {object} response.APIResponse "Erro de validação"
// @Failure 409 {object} response.APIResponse "Email já cadastrado"
// @Router /auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	user, token, err := c.authService.Register(&req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email já cadastrado" {
			status = http.StatusConflict
		}
		response.Error(w, status, err.Error())
		return
	}

	response.Created(w, map[string]interface{}{
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
// @Success 200 {object} response.APIResponse "Login realizado com sucesso"
// @Failure 401 {object} response.APIResponse "Credenciais inválidas"
// @Router /auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	token, err := c.authService.Login(&req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, token)
}

// Refresh renova o access token usando um refresh token válido.
// @Summary Renovar tokens
// @Description Gera novos access + refresh tokens a partir de um refresh token válido.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshRequest true "Refresh token"
// @Success 200 {object} response.APIResponse "Tokens renovados com sucesso"
// @Failure 401 {object} response.APIResponse "Refresh token inválido ou expirado"
// @Router /auth/refresh [post]
func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	token, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, token)
}

// Verify confirma o email do usuário através de um token de verificação.
// @Summary Verificar email
// @Description Marca o email do usuário como verificado usando o token enviado por email.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string false "Token de verificação"
// @Param request body models.VerifyEmailRequest false "Token no body (alternativa)"
// @Success 200 {object} response.APIResponse "Email verificado com sucesso"
// @Failure 400 {object} response.APIResponse "Token inválido ou expirado"
// @Router /auth/verify-email [get]
// @Router /auth/verify-email [post]
func (c *AuthController) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		var req models.VerifyEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
			response.Error(w, http.StatusBadRequest, "token é obrigatório (query param ?token= ou body JSON)")
			return
		}
		token = req.Token
	}

	if err := c.authService.VerifyEmail(token); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Message(w, http.StatusOK, "email verificado com sucesso")
}

// Verify jwt token
// @Summary Verificar token JWT
// @Description Verifica se o token JWT enviado no header Authorization é válido e retorna os claims.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string false "Token de verificação"
// @Param request body models.VerifyEmailRequest false "Token no body (alternativa)"
// @Success 200 {object} response.APIResponse "Token verificado com sucesso"
// @Failure 400 {object} response.APIResponse "Token inválido ou expirado"
// @Router /auth/verify-token [post]
func (c *AuthController) VerifyJWT(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "token de autenticação é obrigatório")
		return
	}

	claims, err := c.authService.VerifyJWT(token)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "token inválido ou expirado")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}
