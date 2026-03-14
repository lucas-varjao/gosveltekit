package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Mode string

const (
	ModeOffset Mode = "offset"
	ModeCursor Mode = "cursor"
)

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

type OffsetQuery struct {
	Page     int
	PageSize int
	Search   string
	Sort     string
	Order    SortDirection
}

type CursorQuery struct {
	PageSize int
	Search   string
	Sort     string
	Order    SortDirection
	After    string
	Before   string
}

type Response[T any] struct {
	Items          []T    `json:"items"`
	Search         string `json:"search,omitempty"`
	Sort           Sort   `json:"sort"`
	PaginationMode Mode   `json:"pagination_mode"`
	Pagination     any    `json:"pagination"`
}

type OffsetMetadata struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type CursorMetadata struct {
	PageSize   int     `json:"page_size"`
	NextCursor *string `json:"next_cursor,omitempty"`
	PrevCursor *string `json:"prev_cursor,omitempty"`
	HasNext    bool    `json:"has_next"`
	HasPrev    bool    `json:"has_prev"`
}

type CursorToken struct {
	Sort      string        `json:"sort"`
	Direction SortDirection `json:"direction"`
	Value     string        `json:"value"`
	ID        uint          `json:"id"`
}

var ErrInvalidCursor = errors.New("cursor inválido")

func EncodeCursor(token CursorToken) (string, error) {
	payload, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecodeCursor(encoded string) (*CursorToken, error) {
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	var token CursorToken
	if err := json.Unmarshal(payload, &token); err != nil {
		return nil, ErrInvalidCursor
	}

	if token.Sort == "" || token.Direction == "" || token.Value == "" || token.ID == 0 {
		return nil, ErrInvalidCursor
	}

	return &token, nil
}
