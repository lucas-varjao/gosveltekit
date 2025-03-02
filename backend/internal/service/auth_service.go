// service/auth_service.go
package service

import (
	"errors"
	"time"

	"gosveltekit/internal/auth"
	"gosveltekit/internal/models"
	"gosveltekit/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserNotActive      = errors.New("usuário inativo")
)

type AuthService struct {
	userRepo     *repository.UserRepository
	tokenService *auth.TokenService
}

func NewAuthService(userRepo *repository.UserRepository, tokenService *auth.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    time.Time   `json:"expires_at"`
	User         models.User `json:"user"`
}

func (s *AuthService) Login(username, password string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.Active {
		return nil, ErrUserNotActive
	}

	// Comparando a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Gerando access token
	accessToken, expiresAt, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Gerando refresh token
	refreshToken, refreshExpiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Atualizando o usuário com o refresh token
	user.RefreshToken = refreshToken
	user.TokenExpiry = refreshExpiresAt
	user.LastLogin = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Removendo dados sensíveis
	user.PasswordHash = ""
	user.RefreshToken = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.Active {
		return nil, ErrUserNotActive
	}

	// Verificando se o token não expirou
	if time.Now().After(user.TokenExpiry) {
		return nil, auth.ErrExpiredToken
	}

	// Gerando novo access token
	accessToken, expiresAt, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Gerar novo refresh token apenas se estiver próximo da expiração
	var newRefreshToken string
	refreshExpiresAt := user.TokenExpiry

	if time.Until(user.TokenExpiry) < 24*time.Hour {
		newRefreshToken, refreshExpiresAt, err = s.tokenService.GenerateRefreshToken()
		if err != nil {
			return nil, err
		}

		user.RefreshToken = newRefreshToken
		user.TokenExpiry = refreshExpiresAt

		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	} else {
		newRefreshToken = refreshToken
	}

	// Removendo dados sensíveis
	user.PasswordHash = ""
	user.RefreshToken = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User:         *user,
	}, nil
}

func (s *AuthService) Logout(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.RefreshToken = ""
	user.TokenExpiry = time.Time{}

	return s.userRepo.Update(user)
}
