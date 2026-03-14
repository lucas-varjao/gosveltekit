package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"gosveltekit/internal/pagination"
	"gosveltekit/internal/service"

	"github.com/gin-gonic/gin"
)

// ListAdminUsers returns a paginated administrative users listing.
func (h *AuthHandler) ListAdminUsers(c *gin.Context) {
	mode := pagination.Mode(c.Query("pagination_mode"))
	if mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrPaginationModeRequired.Error()})
		return
	}

	var input service.ListAdminUsersInput

	switch mode {
	case pagination.ModeOffset:
		offsetInput, err := buildOffsetAdminUsersInput(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
			return
		}
		input = service.ListAdminUsersInput{
			PaginationMode: mode,
			Offset:         offsetInput,
		}
	case pagination.ModeCursor:
		cursorInput, err := buildCursorAdminUsersInput(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
			return
		}
		input = service.ListAdminUsersInput{
			PaginationMode: mode,
			Cursor:         cursorInput,
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrInvalidPaginationMode.Error()})
		return
	}

	users, err := h.authService.ListAdminUsers(input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPaginationModeRequired),
			errors.Is(err, service.ErrInvalidPaginationMode),
			errors.Is(err, service.ErrInvalidAdminUsersQuery),
			errors.Is(err, service.ErrUnsupportedCursorSort):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao listar usuários"})
		}
		return
	}

	c.JSON(http.StatusOK, users)
}

func buildOffsetAdminUsersInput(c *gin.Context) (*pagination.OffsetQuery, error) {
	page, err := parseOptionalPositiveInt(c.Query("page"))
	if err != nil {
		return nil, err
	}

	pageSize, err := parseOptionalPositiveInt(c.Query("page_size"))
	if err != nil {
		return nil, err
	}

	return &pagination.OffsetQuery{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Sort:     c.Query("sort"),
		Order:    pagination.SortDirection(c.Query("order")),
	}, nil
}

func buildCursorAdminUsersInput(c *gin.Context) (*pagination.CursorQuery, error) {
	pageSize, err := parseOptionalPositiveInt(c.Query("page_size"))
	if err != nil {
		return nil, err
	}

	return &pagination.CursorQuery{
		PageSize: pageSize,
		Search:   c.Query("search"),
		Sort:     c.Query("sort"),
		Order:    pagination.SortDirection(c.Query("order")),
		After:    c.Query("after"),
		Before:   c.Query("before"),
	}, nil
}

func parseOptionalPositiveInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}
