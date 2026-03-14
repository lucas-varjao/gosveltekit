package handlers

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"gosveltekit/internal/pagination"

	"github.com/gin-gonic/gin"
)

const (
	defaultMockItemsPage     = 1
	defaultMockItemsPageSize = 8
	maxMockItemsPageSize     = 50
)

type MockPaginationItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Priority  string    `json:"priority"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

var mockPaginationItems = buildMockPaginationItems()

func (h *AuthHandler) ListMockPaginationItems(c *gin.Context) {
	mode := pagination.Mode(c.Query("pagination_mode"))
	if mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pagination_mode é obrigatório"})
		return
	}

	switch mode {
	case pagination.ModeOffset:
		h.listMockPaginationItemsOffset(c)
	case pagination.ModeCursor:
		h.listMockPaginationItemsCursor(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "pagination_mode inválido"})
	}
}

func (h *AuthHandler) listMockPaginationItemsOffset(c *gin.Context) {
	page, err := parseOptionalPositiveInt(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	pageSize, err := parseOptionalPositiveInt(c.Query("page_size"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	if page == 0 {
		page = defaultMockItemsPage
	}
	if pageSize == 0 {
		pageSize = defaultMockItemsPageSize
	}
	if page < 1 || pageSize < 1 || pageSize > maxMockItemsPageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	sortField, direction, ok := normalizeMockSort(c.Query("sort"), pagination.SortDirection(c.Query("order")), true)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	items := filterAndSortMockPaginationItems(c.Query("search"), sortField, direction)
	totalItems := len(items)
	totalPages := max(1, (totalItems+pageSize-1)/pageSize)

	start := (page - 1) * pageSize
	if start > totalItems {
		start = totalItems
	}
	end := min(start+pageSize, totalItems)

	c.JSON(http.StatusOK, pagination.Response[MockPaginationItem]{
		Items:          items[start:end],
		Search:         strings.TrimSpace(c.Query("search")),
		Sort:           pagination.Sort{Field: sortField, Direction: direction},
		PaginationMode: pagination.ModeOffset,
		Pagination: pagination.OffsetMetadata{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int64(totalItems),
			TotalPages: totalPages,
		},
	})
}

func (h *AuthHandler) listMockPaginationItemsCursor(c *gin.Context) {
	pageSize, err := parseOptionalPositiveInt(c.Query("page_size"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}
	if pageSize == 0 {
		pageSize = defaultMockItemsPageSize
	}
	if pageSize < 1 || pageSize > maxMockItemsPageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	after := strings.TrimSpace(c.Query("after"))
	before := strings.TrimSpace(c.Query("before"))
	if after != "" && before != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
		return
	}

	sortField, direction, ok := normalizeMockSort(c.Query("sort"), pagination.SortDirection(c.Query("order")), false)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sort não suportado para paginação cursor"})
		return
	}

	items := filterAndSortMockPaginationItems(c.Query("search"), sortField, direction)

	start := 0
	end := min(pageSize, len(items))
	if after != "" || before != "" {
		token, err := pagination.DecodeCursor(firstNonEmpty(after, before))
		if err != nil || token.Sort != sortField || token.Direction != direction {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
			return
		}

		cursorIndex := findMockCursorIndex(items, token, sortField)
		if cursorIndex == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros de listagem inválidos"})
			return
		}

		if after != "" {
			start = cursorIndex + 1
			end = min(start+pageSize, len(items))
		} else {
			end = cursorIndex
			start = max(0, end-pageSize)
		}
	}

	pageItems := items[start:end]
	hasPrev := start > 0
	hasNext := end < len(items)

	var prevCursor *string
	var nextCursor *string
	if len(pageItems) > 0 {
		if hasPrev {
			cursor, err := encodeMockCursor(pageItems[0], sortField, direction)
			if err == nil {
				prevCursor = &cursor
			}
		}
		if hasNext {
			cursor, err := encodeMockCursor(pageItems[len(pageItems)-1], sortField, direction)
			if err == nil {
				nextCursor = &cursor
			}
		}
	}

	c.JSON(http.StatusOK, pagination.Response[MockPaginationItem]{
		Items:          pageItems,
		Search:         strings.TrimSpace(c.Query("search")),
		Sort:           pagination.Sort{Field: sortField, Direction: direction},
		PaginationMode: pagination.ModeCursor,
		Pagination: pagination.CursorMetadata{
			PageSize:   pageSize,
			NextCursor: nextCursor,
			PrevCursor: prevCursor,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	})
}

func buildMockPaginationItems() []MockPaginationItem {
	categories := []string{"Analytics", "Billing", "Catalog", "Operations", "Support"}
	priorities := []string{"Low", "Medium", "High"}
	items := make([]MockPaginationItem, 0, 48)
	baseDate := time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC)

	for index := 1; index <= 48; index += 1 {
		items = append(items, MockPaginationItem{
			ID:        strconv.Itoa(index),
			Title:     "Mock Item " + strconv.Itoa(index),
			Category:  categories[(index-1)%len(categories)],
			Priority:  priorities[(index-1)%len(priorities)],
			Active:    index%4 != 0,
			CreatedAt: baseDate.Add(time.Duration(index) * 6 * time.Hour),
		})
	}

	return items
}

func filterAndSortMockPaginationItems(
	search string,
	sortField string,
	direction pagination.SortDirection,
) []MockPaginationItem {
	filtered := make([]MockPaginationItem, 0, len(mockPaginationItems))
	term := strings.ToLower(strings.TrimSpace(search))

	for _, item := range mockPaginationItems {
		if term == "" || mockItemMatchesSearch(item, term) {
			filtered = append(filtered, item)
		}
	}

	slices.SortStableFunc(filtered, func(left, right MockPaginationItem) int {
		comparison := compareMockItems(left, right, sortField)
		if comparison == 0 {
			comparison = compareMockItems(left, right, "id")
		}
		if direction == pagination.SortDesc {
			return -comparison
		}
		return comparison
	})

	return filtered
}

func normalizeMockSort(
	sortField string,
	direction pagination.SortDirection,
	allowOffsetOnlyFields bool,
) (string, pagination.SortDirection, bool) {
	normalizedField := strings.TrimSpace(strings.ToLower(sortField))
	if normalizedField == "" {
		normalizedField = "created_at"
	}

	normalizedDirection := pagination.SortDirection(strings.TrimSpace(strings.ToLower(string(direction))))
	if normalizedDirection == "" {
		normalizedDirection = pagination.SortDesc
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"title":      true,
		"category":   true,
	}
	if allowOffsetOnlyFields {
		allowedSorts["priority"] = true
	}

	if !allowedSorts[normalizedField] {
		return "", "", false
	}
	if normalizedDirection != pagination.SortAsc && normalizedDirection != pagination.SortDesc {
		return "", "", false
	}

	return normalizedField, normalizedDirection, true
}

func mockItemMatchesSearch(item MockPaginationItem, term string) bool {
	return strings.Contains(strings.ToLower(item.Title), term) ||
		strings.Contains(strings.ToLower(item.Category), term) ||
		strings.Contains(strings.ToLower(item.Priority), term)
}

func compareMockItems(left MockPaginationItem, right MockPaginationItem, sortField string) int {
	switch sortField {
	case "title":
		return strings.Compare(left.Title, right.Title)
	case "category":
		return strings.Compare(left.Category, right.Category)
	case "priority":
		return strings.Compare(left.Priority, right.Priority)
	case "created_at":
		return compareTimes(left.CreatedAt, right.CreatedAt)
	case "id":
		leftID, leftErr := strconv.Atoi(left.ID)
		rightID, rightErr := strconv.Atoi(right.ID)
		if leftErr == nil && rightErr == nil {
			switch {
			case leftID < rightID:
				return -1
			case leftID > rightID:
				return 1
			default:
				return 0
			}
		}
		return strings.Compare(left.ID, right.ID)
	default:
		return 0
	}
}

func compareTimes(left time.Time, right time.Time) int {
	switch {
	case left.Before(right):
		return -1
	case left.After(right):
		return 1
	default:
		return 0
	}
}

func findMockCursorIndex(
	items []MockPaginationItem,
	token *pagination.CursorToken,
	sortField string,
) int {
	for index, item := range items {
		if item.ID != strconv.FormatUint(uint64(token.ID), 10) {
			continue
		}
		if mockCursorValue(item, sortField) == token.Value {
			return index
		}
	}

	return -1
}

func encodeMockCursor(
	item MockPaginationItem,
	sortField string,
	direction pagination.SortDirection,
) (string, error) {
	id, err := strconv.ParseUint(item.ID, 10, 64)
	if err != nil {
		return "", err
	}

	return pagination.EncodeCursor(pagination.CursorToken{
		Sort:      sortField,
		Direction: direction,
		Value:     mockCursorValue(item, sortField),
		ID:        uint(id),
	})
}

func mockCursorValue(item MockPaginationItem, sortField string) string {
	switch sortField {
	case "title":
		return item.Title
	case "category":
		return item.Category
	default:
		return item.CreatedAt.UTC().Format(time.RFC3339Nano)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
