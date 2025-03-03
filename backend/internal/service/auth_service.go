// service/auth_service.go
package service

import (
	"errors"
	"sync"
	"time"

	"gosveltekit/internal/auth"
	"gosveltekit/internal/models"
	"gosveltekit/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserNotActive      = errors.New("usuário inativo")
	ErrInvalidToken       = errors.New("token inválido")
	ErrExpiredToken       = errors.New("token expirado")
)

type AuthService struct {
	userRepo            *repository.UserRepository
	tokenService        *auth.TokenService
	failedLoginAttempts map[string]int
	failedLoginMutex    sync.RWMutex
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

func (s *AuthService) Login(username, password string, ip, userAgent string) (*LoginResponse, error) {
	// Verificar tentativas falhas
	s.failedLoginMutex.RLock()
	attempts, exists := s.failedLoginAttempts[username]
	s.failedLoginMutex.RUnlock()

	if exists && attempts >= 5 {
		return nil, errors.New("conta temporariamente bloqueada, tente novamente mais tarde")
	}
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		s.incrementFailedLoginAttempt(username)
		return nil, ErrInvalidCredentials
	}

	if !user.Active {
		s.incrementFailedLoginAttempt(username)
		return nil, ErrUserNotActive
	}

	// Comparando a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.incrementFailedLoginAttempt(username)
		return nil, ErrInvalidCredentials
	}

	// Gerando access token
	accessToken, expiresAt, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		s.incrementFailedLoginAttempt(username)
		return nil, err
	}

	// Gerando refresh token
	refreshToken, refreshExpiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		s.incrementFailedLoginAttempt(username)
		return nil, err
	}

	// Atualizando o usuário com o refresh token
	user.RefreshToken = refreshToken
	user.TokenExpiry = refreshExpiresAt
	user.LastLogin = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		s.incrementFailedLoginAttempt(username)
		return nil, err
	}

	s.resetFailedLoginAttempts(username)

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

	if time.Until(user.TokenExpiry) < 24*time.Hour {
		var tokenExpiry time.Time
		newRefreshToken, tokenExpiry, err = s.tokenService.GenerateRefreshToken()
		if err != nil {
			return nil, err
		}

		user.RefreshToken = newRefreshToken
		user.TokenExpiry = tokenExpiry

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

func (s *AuthService) Logout(userID uint, accessToken string) error {
	// Adicionamos o token à blacklist
	if err := s.tokenService.BlacklistToken(accessToken); err != nil {
		return err
	}

	// Buscamos o usuário para limpar o refresh token
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Limpamos o refresh token
	user.RefreshToken = ""
	user.TokenExpiry = time.Time{}

	return s.userRepo.Update(user)
}

func (s *AuthService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Não informamos se o email existe ou não por segurança
		return nil
	}

	// Gerar token de recuperação (válido por 1 hora)
	resetToken, expiresAt, err := s.tokenService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	// Armazenar token no usuário
	user.ResetToken = resetToken
	user.ResetTokenExpiry = expiresAt
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Enviar email com link de recuperação
	// resetLink := fmt.Sprintf("https://seu-app.com/reset-password?token=%s", resetToken)
	// emailService.SendPasswordResetEmail(user.Email, resetLink)

	return nil
}

func (s *AuthService) ResetPassword(resetToken, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(resetToken)
	if err != nil {
		return ErrInvalidToken
	}

	// Verificar se o token ainda é válido
	if time.Now().After(user.ResetTokenExpiry) {
		return ErrExpiredToken
	}

	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Atualizar senha e limpar tokens
	user.PasswordHash = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}

	return s.userRepo.Update(user)
}

func (s *AuthService) incrementFailedLoginAttempt(username string) {
	s.failedLoginMutex.Lock()
	defer s.failedLoginMutex.Unlock()

	if _, exists := s.failedLoginAttempts[username]; !exists {
		s.failedLoginAttempts[username] = 0
	}
	s.failedLoginAttempts[username]++

	// Agendar limpeza após 30 minutos
	go func(username string) {
		time.Sleep(30 * time.Minute)
		s.resetFailedLoginAttempts(username)
	}(username)
}

func (s *AuthService) resetFailedLoginAttempts(username string) {
	s.failedLoginMutex.Lock()
	defer s.failedLoginMutex.Unlock()

	delete(s.failedLoginAttempts, username)
}
