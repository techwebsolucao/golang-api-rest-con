package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/user/golang-api-rest/internal/config"
)

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int, role string, cfg *config.Config) (string, time.Time, error) {
	expiresAt := time.Now().Add(15 * time.Minute)

	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func GenerateRefreshToken(userID int, role string, cfg *config.Config) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateToken(tokenStr string, cfg *config.Config) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}

func GenerateVerificationToken(userID int, cfg *config.Config) (string, error) {
	expiresAt := time.Now().Add(1 * time.Hour)

	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}