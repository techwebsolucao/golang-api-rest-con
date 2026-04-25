package services

import (
	"errors"
	"fmt"

	"github.com/user/golang-api-rest/internal/config"
	"github.com/user/golang-api-rest/internal/models"
	"github.com/user/golang-api-rest/internal/repositories"
	"github.com/user/golang-api-rest/internal/utils"
)

type AuthService struct {
	repo  repositories.UserRepository
	cfg   *config.Config
	email *EmailService
}

func NewAuthService(repo repositories.UserRepository, cfg *config.Config, email *EmailService) *AuthService {
	return &AuthService{repo: repo, cfg: cfg, email: email}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, *models.TokenResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, nil, errors.New("nome, email e senha são obrigatórios")
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, nil, err
	}

	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, nil, errors.New("email já cadastrado")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "user",
		Verified:     false,
	}

	id, err := s.repo.Create(user)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}
	user.ID = id

	verifyToken, err := utils.GenerateVerificationToken(user.ID, s.cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao gerar token de verificação: %w", err)
	}

	go s.email.SendVerificationEmail(user.Email, user.Name, verifyToken)

	accessToken, expiresAt, err := utils.GenerateAccessToken(user.ID, user.Role, s.cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role, s.cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	tokenResp := &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	return user, tokenResp, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.TokenResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email e senha são obrigatórios")
	}

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("credenciais inválidas")
	}

	if !user.Verified {
		return nil, errors.New("email não verificado. Por favor, verifique seu email antes de fazer login.")
	}

	accessToken, expiresAt, err := utils.GenerateAccessToken(user.ID, user.Role, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, s.cfg)
	if err != nil {
		return nil, errors.New("refresh token inválido ou expirado")
	}

	accessToken, expiresAt, err := utils.GenerateAccessToken(claims.UserID, claims.Role, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	newRefresh, err := utils.GenerateRefreshToken(claims.UserID, claims.Role, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *AuthService) VerifyEmail(tokenStr string) error {
	claims, err := utils.ValidateToken(tokenStr, s.cfg)
	if err != nil {
		return errors.New("token de verificação inválido ou expirado")
	}

	user, err := s.repo.GetByID(claims.UserID)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	if user.Verified {
		return errors.New("email já verificado")
	}

	user.Verified = true
	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	go s.email.SendWelcomeEmail(user.Email, user.Name)
	return nil
}

func (s *AuthService) VerifyJWT(tokenStr string) (*utils.TokenClaims, error) {
	claims, err := utils.ValidateToken(tokenStr, s.cfg)
	if err != nil {
		return nil, errors.New("token inválido ou expirado")
	}
	return claims, nil
}
