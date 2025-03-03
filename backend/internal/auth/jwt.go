// auth/jwt.go
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
	config    *config.Config
	blacklist *TokenBlacklist
}

func NewTokenService(cfg *config.Config) *TokenService {
	blacklist := NewTokenBlacklist(1 * time.Hour)
	return &TokenService{config: cfg, blacklist: blacklist}
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
	if ts.blacklist.IsBlacklisted(tokenString) {
		return nil, ErrInvalidToken
	}
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

// GeneratePasswordResetToken creates a token for password reset
// Now returns both the plaintext token (for the email link) and the hashed token (for storage)
func (ts *TokenService) GeneratePasswordResetToken(userID uint) (string, string, time.Time, error) {
	expiresAt := time.Now().Add(ts.config.JWT.PasswordResetTTL)

	// Generate a secure random token
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Create the plaintext token to be sent to the user
	plaintextToken := base64.URLEncoding.EncodeToString(randomBytes)

	// Create a hash of the token to store in the database
	hashedToken := hashToken(plaintextToken)

	// Create JWT claims for additional security (optional but keeps consistency)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    ts.config.JWT.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			// Adding a special purpose for the token
			ID: "password-reset",
		},
	}

	// Sign the token with our secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ts.config.JWT.SecretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// For password reset, we return the plaintext token for the email link
	// and the hashed version for storage
	return plaintextToken + "." + tokenString, hashedToken, expiresAt, nil
}

// hashToken creates a SHA-256 hash of the token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// VerifyPasswordResetToken verifies if a password reset token is valid
func (ts *TokenService) VerifyPasswordResetToken(plaintextToken, hashedTokenFromDB string) bool {
	// Hash the plaintext token and compare with the stored hash
	computedHash := hashToken(plaintextToken)
	return computedHash == hashedTokenFromDB
}

func (ts *TokenService) BlacklistToken(tokenString string) error {
	// Primeiro, validamos o token para obter seu tempo de expiração
	claims, err := ts.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// Adicionamos à blacklist com o tempo de expiração
	expTime := claims.ExpiresAt.Time
	ts.blacklist.Add(tokenString, expTime)

	return nil
}

func (ts *TokenService) IsTokenBlacklisted(tokenString string) bool {
	return ts.blacklist.IsBlacklisted(tokenString)
}
