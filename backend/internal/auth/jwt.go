// auth/jwt.go
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"gosveltekit/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
)

// Claims personalizado para o payload do JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	config *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{config: cfg}
}

// GenerateAccessToken gera um token JWT
func (ts *TokenService) GenerateAccessToken(userID uint, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(ts.config.JWT.AccessTokenTTL)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    ts.config.JWT.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(ts.config.JWT.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken gera um token de atualização
func (ts *TokenService) GenerateRefreshToken() (string, time.Time, error) {
	expiresAt := time.Now().Add(ts.config.JWT.RefreshTokenTTL)

	// Gerando um token aleatório (você pode usar uuid ou outro método)
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", time.Time{}, err
	}

	tokenString := base64.URLEncoding.EncodeToString(randomBytes)
	return tokenString, expiresAt, nil
}

// ValidateToken valida um token JWT
func (ts *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Verifica o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(ts.config.JWT.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
