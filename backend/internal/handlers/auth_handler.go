// handlers/auth_handler.go
package handlers

import (
	"net/http"

	"gosveltekit/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get client IP and user agent
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	response, err := h.authService.Login(req.Username, req.Password, ip, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		message := "credenciais inválidas"

		if err == service.ErrUserNotActive {
			message = "usuário inativo"
		} else if err.Error() == "conta temporariamente bloqueada, tente novamente mais tarde" {
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		message := "token de atualização inválido"

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
		return
	}

	// Obtém o token atual do contexto
	accessToken, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token não encontrado"})
		return
	}

	if err := h.authService.Logout(userID.(uint), accessToken.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao fazer logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout realizado com sucesso"})
}
