package gorm

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"gosveltekit/internal/models"
	"gosveltekit/internal/pagination"
)

var adminUsersSortColumns = map[string]string{
	"created_at":   "created_at",
	"display_name": "display_name",
	"email":        "email",
	"identifier":   "username",
	"last_login":   "last_login",
	"role":         "role",
}

type CursorPageResult struct {
	Users      []*models.User
	NextCursor *string
	PrevCursor *string
	HasNext    bool
	HasPrev    bool
}

// ListUsersOffset returns offset-based paginated users for administrative tables.
func (a *UserAdapter) ListUsersOffset(
	input pagination.OffsetQuery,
) ([]*models.User, int64, error) {
	query := a.db.Model(&models.User{})

	search := strings.TrimSpace(input.Search)
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

	orderColumn, ok := adminUsersSortColumns[input.Sort]
	if !ok {
		orderColumn = "created_at"
	}

	var users []*models.User
	if err := query.
		Order(orderColumn + " " + string(input.Order)).
		Limit(input.PageSize).
		Offset((input.Page - 1) * input.PageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalItems, nil
}

// ListUsersCursor returns cursor-based paginated users for administrative tables.
func (a *UserAdapter) ListUsersCursor(input pagination.CursorQuery) (*CursorPageResult, error) {
	query := a.db.Model(&models.User{})

	search := strings.TrimSpace(input.Search)
	if search != "" {
		likeTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(display_name) LIKE ?",
			likeTerm,
			likeTerm,
			likeTerm,
		)
	}

	orderColumn, ok := adminUsersSortColumns[input.Sort]
	if !ok {
		orderColumn = "created_at"
	}

	var (
		users      []*models.User
		hasNext    bool
		hasPrev    bool
		nextCursor *string
		prevCursor *string
	)

	limit := input.PageSize + 1
	isBefore := input.Before != ""
	cursorValue := input.After
	if isBefore {
		cursorValue = input.Before
	}

	if cursorValue != "" {
		token, err := pagination.DecodeCursor(cursorValue)
		if err != nil {
			return nil, err
		}

		if token.Sort != input.Sort || token.Direction != input.Order {
			return nil, pagination.ErrInvalidCursor
		}

		cursorQuery, reverseOrder, err := buildCursorQuery(orderColumn, input.Order, token, isBefore)
		if err != nil {
			return nil, err
		}

		query = query.Where(cursorQuery.sql, cursorQuery.args...)
		if err := query.
			Order(fmt.Sprintf("%s %s", orderColumn, reverseOrder.primary)).
			Order(fmt.Sprintf("id %s", reverseOrder.tieBreaker)).
			Limit(limit).
			Find(&users).Error; err != nil {
			return nil, err
		}

		if len(users) > input.PageSize {
			if isBefore {
				hasPrev = true
			} else {
				hasNext = true
			}
			users = users[:input.PageSize]
		}

		if isBefore {
			slices.Reverse(users)
			hasNext = true
		} else {
			hasPrev = true
		}
	} else {
		order := buildOrder(input.Order, false)
		if err := query.
			Order(fmt.Sprintf("%s %s", orderColumn, order.primary)).
			Order(fmt.Sprintf("id %s", order.tieBreaker)).
			Limit(limit).
			Find(&users).Error; err != nil {
			return nil, err
		}

		if len(users) > input.PageSize {
			hasNext = true
			users = users[:input.PageSize]
		}
	}

	if len(users) > 0 {
		if hasPrev {
			cursor, err := encodeUserCursor(users[0], input.Sort, input.Order)
			if err != nil {
				return nil, err
			}
			prevCursor = &cursor
		}
		if hasNext {
			cursor, err := encodeUserCursor(users[len(users)-1], input.Sort, input.Order)
			if err != nil {
				return nil, err
			}
			nextCursor = &cursor
		}
	}

	return &CursorPageResult{
		Users:      users,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}, nil
}

type cursorWhereClause struct {
	sql  string
	args []any
}

type orderParts struct {
	primary    string
	tieBreaker string
}

func buildCursorQuery(
	orderColumn string,
	direction pagination.SortDirection,
	token *pagination.CursorToken,
	isBefore bool,
) (*cursorWhereClause, orderParts, error) {
	value, err := parseCursorValue(token.Sort, token.Value)
	if err != nil {
		return nil, orderParts{}, pagination.ErrInvalidCursor
	}

	order := buildOrder(direction, isBefore)
	comparator := ">"
	tieComparator := ">"
	if direction == pagination.SortDesc {
		comparator = "<"
		tieComparator = "<"
	}
	if isBefore {
		if comparator == ">" {
			comparator = "<"
			tieComparator = "<"
		} else {
			comparator = ">"
			tieComparator = ">"
		}
	}

	return &cursorWhereClause{
		sql: fmt.Sprintf("(%s %s ?) OR (%s = ? AND id %s ?)", orderColumn, comparator, orderColumn, tieComparator),
		args: []any{
			value,
			value,
			token.ID,
		},
	}, order, nil
}

func buildOrder(direction pagination.SortDirection, reverse bool) orderParts {
	primary := string(direction)
	tieBreaker := string(direction)
	if reverse {
		if direction == pagination.SortAsc {
			primary = string(pagination.SortDesc)
			tieBreaker = string(pagination.SortDesc)
		} else {
			primary = string(pagination.SortAsc)
			tieBreaker = string(pagination.SortAsc)
		}
	}

	return orderParts{
		primary:    primary,
		tieBreaker: tieBreaker,
	}
}

func parseCursorValue(sortField string, raw string) (any, error) {
	switch sortField {
	case "created_at":
		return time.Parse(time.RFC3339Nano, raw)
	default:
		return raw, nil
	}
}

func encodeUserCursor(
	user *models.User,
	sortField string,
	direction pagination.SortDirection,
) (string, error) {
	value := ""
	switch sortField {
	case "created_at":
		value = user.CreatedAt.UTC().Format(time.RFC3339Nano)
	case "email":
		value = user.Email
	case "identifier":
		value = user.Username
	default:
		return "", pagination.ErrInvalidCursor
	}

	return pagination.EncodeCursor(pagination.CursorToken{
		Sort:      sortField,
		Direction: direction,
		Value:     value,
		ID:        user.ID,
	})
}
