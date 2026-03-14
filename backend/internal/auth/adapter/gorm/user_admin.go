package gorm

import (
	"strings"

	"gosveltekit/internal/models"
)

var adminUsersSortColumns = map[string]string{
	"created_at":   "created_at",
	"display_name": "display_name",
	"email":        "email",
	"identifier":   "username",
	"last_login":   "last_login",
	"role":         "role",
}

// ListUsers returns paginated users for administrative tables.
func (a *UserAdapter) ListUsers(
	page int,
	pageSize int,
	search string,
	sortField string,
	sortOrder string,
) ([]*models.User, int64, error) {
	query := a.db.Model(&models.User{})

	search = strings.TrimSpace(search)
	if search != "" {
		likeTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(display_name) LIKE ?",
			likeTerm,
			likeTerm,
			likeTerm,
		)
	}

	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	orderColumn, ok := adminUsersSortColumns[sortField]
	if !ok {
		orderColumn = "created_at"
	}

	var users []*models.User
	if err := query.
		Order(orderColumn + " " + sortOrder).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalItems, nil
}
