package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gosveltekit/internal/models"
)

const (
	defaultAdminUsersPage     = 1
	defaultAdminUsersPageSize = 10
	maxAdminUsersPageSize     = 100
)

var ErrInvalidAdminUsersQuery = errors.New("parâmetros de listagem inválidos")

type ListAdminUsersInput struct {
	Page     int
	PageSize int
	Search   string
	Sort     string
	Order    string
}

type AdminUsersSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type PaginatedResult[T any] struct {
	Items      []T            `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalItems int64          `json:"total_items"`
	TotalPages int            `json:"total_pages"`
	Sort       AdminUsersSort `json:"sort"`
	Search     string         `json:"search,omitempty"`
}

type AdminUserRow struct {
	ID          string    `json:"id"`
	Identifier  string    `json:"identifier"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Role        string    `json:"role"`
	Active      bool      `json:"active"`
	LastLogin   time.Time `json:"last_login"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *AuthService) ListAdminUsers(input ListAdminUsersInput) (*PaginatedResult[AdminUserRow], error) {
	normalized, err := normalizeAdminUsersInput(input)
	if err != nil {
		return nil, err
	}

	users, totalItems, err := s.userAdapter.ListUsers(
		normalized.Page,
		normalized.PageSize,
		normalized.Search,
		normalized.Sort,
		normalized.Order,
	)
	if err != nil {
		return nil, err
	}

	items := make([]AdminUserRow, 0, len(users))
	for _, user := range users {
		items = append(items, toAdminUserRow(user))
	}

	totalPages := int((totalItems + int64(normalized.PageSize) - 1) / int64(normalized.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedResult[AdminUserRow]{
		Items:      items,
		Page:       normalized.Page,
		PageSize:   normalized.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Sort: AdminUsersSort{
			Field:     normalized.Sort,
			Direction: normalized.Order,
		},
		Search: normalized.Search,
	}, nil
}

func normalizeAdminUsersInput(input ListAdminUsersInput) (ListAdminUsersInput, error) {
	page := input.Page
	if page == 0 {
		page = defaultAdminUsersPage
	}
	if page < 1 {
		return ListAdminUsersInput{}, ErrInvalidAdminUsersQuery
	}

	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = defaultAdminUsersPageSize
	}
	if pageSize < 1 || pageSize > maxAdminUsersPageSize {
		return ListAdminUsersInput{}, ErrInvalidAdminUsersQuery
	}

	sortField := strings.TrimSpace(strings.ToLower(input.Sort))
	if sortField == "" {
		sortField = "created_at"
	}

	order := strings.TrimSpace(strings.ToLower(input.Order))
	if order == "" {
		order = "desc"
	}

	allowedSorts := map[string]bool{
		"created_at":   true,
		"display_name": true,
		"email":        true,
		"identifier":   true,
		"last_login":   true,
		"role":         true,
	}

	if !allowedSorts[sortField] {
		return ListAdminUsersInput{}, ErrInvalidAdminUsersQuery
	}

	if order != "asc" && order != "desc" {
		return ListAdminUsersInput{}, ErrInvalidAdminUsersQuery
	}

	return ListAdminUsersInput{
		Page:     page,
		PageSize: pageSize,
		Search:   strings.TrimSpace(input.Search),
		Sort:     sortField,
		Order:    order,
	}, nil
}

func toAdminUserRow(user *models.User) AdminUserRow {
	return AdminUserRow{
		ID:          strconv.FormatUint(uint64(user.ID), 10),
		Identifier:  user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		Active:      user.Active,
		LastLogin:   user.LastLogin,
		CreatedAt:   user.CreatedAt,
	}
}
