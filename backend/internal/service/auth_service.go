// backend/internal/service/auth_service.go

package service

import (
	"errors"
	"strings"
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

// AuthServiceInterface defines the methods that an auth service must implement
type AuthServiceInterface interface {
	Login(username, password, ip, userAgent string) (*LoginResponse, error)
	RefreshToken(refreshToken string) (*LoginResponse, error)
	Logout(userID uint, accessToken string) error
	Register(username, email, password, displayName string) (*models.User, error)
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error
}

type AuthService struct {
	userRepo            *repository.UserRepository
	tokenService        *auth.TokenService
	failedLoginAttempts map[string]int
	failedLoginMutex    sync.RWMutex
}

func NewAuthService(userRepo *repository.UserRepository, tokenService *auth.TokenService) *AuthService {
	return &AuthService{
		userRepo:            userRepo,
		tokenService:        tokenService,
		failedLoginAttempts: make(map[string]int),
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
	user.AccessTokenExpiry = expiresAt
	user.RefreshTokenExpiry = refreshExpiresAt
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
	if time.Now().After(user.RefreshTokenExpiry) {
		return nil, auth.ErrExpiredToken
	}

	// Gerando novo access token
	accessToken, expiresAt, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	user.AccessTokenExpiry = expiresAt

	// Gerar novo refresh token apenas se estiver próximo da expiração
	var newRefreshToken string

	if time.Until(user.RefreshTokenExpiry) < 24*time.Hour {
		var tokenExpiry time.Time
		newRefreshToken, tokenExpiry, err = s.tokenService.GenerateRefreshToken()
		if err != nil {
			return nil, err
		}

		user.RefreshToken = newRefreshToken
		user.RefreshTokenExpiry = tokenExpiry

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
	user.AccessTokenExpiry = time.Time{}
	user.RefreshTokenExpiry = time.Time{}

	return s.userRepo.Update(user)
}

func (s *AuthService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Não informamos se o email existe ou não por segurança
		return nil
	}

	// Explicitly invalidate any existing reset tokens for this user first
	// by setting it to an empty string and resetting the expiry
	if user.ResetToken != "" {
		user.ResetToken = ""
		user.ResetTokenExpiry = time.Time{}
		if err := s.userRepo.Update(user); err != nil {
			return err
		}
	}

	// Gerar token de recuperação (válido por 1 hora)
	// Now using the enhanced token generation that returns both plaintext and hashed tokens
	plaintextToken, hashedToken, expiresAt, err := s.tokenService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	// Armazenar token HASH no usuário (não o token original)
	user.ResetToken = hashedToken
	user.ResetTokenExpiry = expiresAt
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Enviar email com link de recuperação
	// O token plaintext é usado para o link de recuperação
	// Temporarily log the token for development purposes
	_ = plaintextToken // TODO: Implement email service to send reset link
	// resetLink := fmt.Sprintf("https://seu-app.com/reset-password?token=%s", plaintextToken)
	// emailService.SendPasswordResetEmail(user.Email, resetLink)

	return nil
}

func (s *AuthService) ResetPassword(tokenFromUser, newPassword string) error {
	// Split the token to get the plaintext part
	parts := strings.SplitN(tokenFromUser, ".", 2)
	if len(parts) != 2 {
		return ErrInvalidToken
	}
	plaintextToken := parts[0]

	// Find all users with non-expired reset tokens
	users, err := s.userRepo.FindUsersWithResetTokens()
	if err != nil {
		return err
	}

	// Find the user with the matching hashed token
	var matchedUser *models.User
	for _, user := range users {
		// Verify if current time is before token expiry
		if time.Now().After(user.ResetTokenExpiry) {
			continue // Skip expired tokens
		}

		// Check if the plaintext token hash matches the stored hash
		if s.tokenService.VerifyPasswordResetToken(plaintextToken, user.ResetToken) {
			matchedUser = user
			break
		}
	}

	if matchedUser == nil {
		return ErrInvalidToken
	}

	// Verificar se o token ainda é válido
	if time.Now().After(matchedUser.ResetTokenExpiry) {
		return ErrExpiredToken
	}

	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Atualizar senha e limpar tokens
	matchedUser.PasswordHash = string(hashedPassword)
	matchedUser.ResetToken = ""
	matchedUser.ResetTokenExpiry = time.Time{}

	return s.userRepo.Update(matchedUser)
}

func (s *AuthService) Register(username, email, password, displayName string) (*models.User, error) {
	// Check if the user already exists by username
	userByUsername, _ := s.userRepo.FindByUsername(username)
	if userByUsername != nil {
		return nil, errors.New("username already exists")
	}

	// Check if the user already exists by email
	userByEmail, _ := s.userRepo.FindByEmail(email)
	if userByEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		DisplayName:  displayName,
	}

	// Save the user to the database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
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
