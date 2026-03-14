package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"gosveltekit/internal/service"

	"github.com/gin-gonic/gin"
)

// ListAdminUsers returns a paginated administrative users listing.
func (h *AuthHandler) ListAdminUsers(c *gin.Context) {
	page, err := parsePositiveInt(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	pageSize, err := parsePositiveInt(c.DefaultQuery("page_size", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	users, err := h.authService.ListAdminUsers(service.ListAdminUsersInput{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Sort:     c.Query("sort"),
		Order:    c.Query("order"),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAdminUsersQuery):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao listar usuários"})
		}
		return
	}

	c.JSON(http.StatusOK, users)
}

func parsePositiveInt(value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}
