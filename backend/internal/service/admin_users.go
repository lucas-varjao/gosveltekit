package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gosveltekit/internal/models"
	"gosveltekit/internal/pagination"
)

const (
	defaultAdminUsersPage     = 1
	defaultAdminUsersPageSize = 10
	maxAdminUsersPageSize     = 100
)

var (
	ErrPaginationModeRequired = errors.New("pagination_mode é obrigatório")
	ErrInvalidPaginationMode  = errors.New("pagination_mode inválido")
	ErrInvalidAdminUsersQuery = errors.New("parâmetros de listagem inválidos")
	ErrUnsupportedCursorSort  = errors.New("sort não suportado para paginação cursor")
)

type ListAdminUsersInput struct {
	PaginationMode pagination.Mode
	Offset         *pagination.OffsetQuery
	Cursor         *pagination.CursorQuery
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

func (s *AuthService) ListAdminUsers(input ListAdminUsersInput) (*pagination.Response[AdminUserRow], error) {
	switch input.PaginationMode {
	case "":
		return nil, ErrPaginationModeRequired
	case pagination.ModeOffset:
		return s.listAdminUsersOffset(input.Offset)
	case pagination.ModeCursor:
		return s.listAdminUsersCursor(input.Cursor)
	default:
		return nil, ErrInvalidPaginationMode
	}
}

func (s *AuthService) listAdminUsersOffset(
	input *pagination.OffsetQuery,
) (*pagination.Response[AdminUserRow], error) {
	normalized, err := normalizeAdminUsersOffsetInput(input)
	if err != nil {
		return nil, err
	}

	users, totalItems, err := s.userAdapter.ListUsersOffset(normalized)
	if err != nil {
		return nil, err
	}

	totalPages := int((totalItems + int64(normalized.PageSize) - 1) / int64(normalized.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &pagination.Response[AdminUserRow]{
		Items:          toAdminUserRows(users),
		Search:         normalized.Search,
		Sort:           pagination.Sort{Field: normalized.Sort, Direction: normalized.Order},
		PaginationMode: pagination.ModeOffset,
		Pagination: pagination.OffsetMetadata{
			Page:       normalized.Page,
			PageSize:   normalized.PageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *AuthService) listAdminUsersCursor(
	input *pagination.CursorQuery,
) (*pagination.Response[AdminUserRow], error) {
	normalized, err := normalizeAdminUsersCursorInput(input)
	if err != nil {
		return nil, err
	}

	result, err := s.userAdapter.ListUsersCursor(normalized)
	if err != nil {
		if errors.Is(err, pagination.ErrInvalidCursor) {
			return nil, ErrInvalidAdminUsersQuery
		}
		return nil, err
	}

	return &pagination.Response[AdminUserRow]{
		Items:          toAdminUserRows(result.Users),
		Search:         normalized.Search,
		Sort:           pagination.Sort{Field: normalized.Sort, Direction: normalized.Order},
		PaginationMode: pagination.ModeCursor,
		Pagination: pagination.CursorMetadata{
			PageSize:   normalized.PageSize,
			NextCursor: result.NextCursor,
			PrevCursor: result.PrevCursor,
			HasNext:    result.HasNext,
			HasPrev:    result.HasPrev,
		},
	}, nil
}

func normalizeAdminUsersOffsetInput(input *pagination.OffsetQuery) (pagination.OffsetQuery, error) {
	if input == nil {
		input = &pagination.OffsetQuery{}
	}

	page := input.Page
	if page == 0 {
		page = defaultAdminUsersPage
	}
	if page < 1 {
		return pagination.OffsetQuery{}, ErrInvalidAdminUsersQuery
	}

	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = defaultAdminUsersPageSize
	}
	if pageSize < 1 || pageSize > maxAdminUsersPageSize {
		return pagination.OffsetQuery{}, ErrInvalidAdminUsersQuery
	}

	sortField, order, err := normalizeAdminUsersSort(input.Sort, input.Order)
	if err != nil {
		return pagination.OffsetQuery{}, err
	}

	return pagination.OffsetQuery{
		Page:     page,
		PageSize: pageSize,
		Search:   strings.TrimSpace(input.Search),
		Sort:     sortField,
		Order:    order,
	}, nil
}

func normalizeAdminUsersCursorInput(input *pagination.CursorQuery) (pagination.CursorQuery, error) {
	if input == nil {
		input = &pagination.CursorQuery{}
	}

	if strings.TrimSpace(input.After) != "" && strings.TrimSpace(input.Before) != "" {
		return pagination.CursorQuery{}, ErrInvalidAdminUsersQuery
	}

	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = defaultAdminUsersPageSize
	}
	if pageSize < 1 || pageSize > maxAdminUsersPageSize {
		return pagination.CursorQuery{}, ErrInvalidAdminUsersQuery
	}

	sortField, order, err := normalizeAdminUsersSort(input.Sort, input.Order)
	if err != nil {
		return pagination.CursorQuery{}, err
	}

	allowedCursorSorts := map[string]bool{
		"created_at": true,
		"email":      true,
		"identifier": true,
	}
	if !allowedCursorSorts[sortField] {
		return pagination.CursorQuery{}, ErrUnsupportedCursorSort
	}

	return pagination.CursorQuery{
		PageSize: pageSize,
		Search:   strings.TrimSpace(input.Search),
		Sort:     sortField,
		Order:    order,
		After:    strings.TrimSpace(input.After),
		Before:   strings.TrimSpace(input.Before),
	}, nil
}

func normalizeAdminUsersSort(sortField string, order pagination.SortDirection) (string, pagination.SortDirection, error) {
	normalizedField := strings.TrimSpace(strings.ToLower(sortField))
	if normalizedField == "" {
		normalizedField = "created_at"
	}

	normalizedOrder := pagination.SortDirection(strings.TrimSpace(strings.ToLower(string(order))))
	if normalizedOrder == "" {
		normalizedOrder = pagination.SortDesc
	}

	allowedSorts := map[string]bool{
		"created_at":   true,
		"display_name": true,
		"email":        true,
		"identifier":   true,
		"last_login":   true,
		"role":         true,
	}
	if !allowedSorts[normalizedField] {
		return "", "", ErrInvalidAdminUsersQuery
	}

	if normalizedOrder != pagination.SortAsc && normalizedOrder != pagination.SortDesc {
		return "", "", ErrInvalidAdminUsersQuery
	}

	return normalizedField, normalizedOrder, nil
}

func toAdminUserRows(users []*models.User) []AdminUserRow {
	items := make([]AdminUserRow, 0, len(users))
	for _, user := range users {
		items = append(items, toAdminUserRow(user))
	}
	return items
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
