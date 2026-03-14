package gorm

import (
	"testing"
	"time"

	"gosveltekit/internal/models"
	"gosveltekit/internal/pagination"
	"gosveltekit/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupAdminUsersTestDB(t *testing.T) (*gorm.DB, *UserAdapter) {
	t.Helper()

	db := testutil.NewSQLiteTestDB(t, &models.User{})
	return db, NewUserAdapter(db)
}

func seedAdminUsers(t *testing.T, db *gorm.DB, users []*models.User) {
	t.Helper()

	for _, user := range users {
		require.NoError(t, db.Create(user).Error)
	}
}

func TestUserAdapter_ListUsersOffset(t *testing.T) {
	db, adapter := setupAdminUsersTestDB(t)
	seedAdminUsers(t, db, []*models.User{
		{Username: "zoe", Email: "zoe@example.com", DisplayName: "Zoe Doe", PasswordHash: "hash", Role: "user"},
		{Username: "anna", Email: "anna@example.com", DisplayName: "Anna Doe", PasswordHash: "hash", Role: "manager"},
	})

	users, totalItems, err := adapter.ListUsersOffset(pagination.OffsetQuery{
		Page:     1,
		PageSize: 1,
		Search:   "doe",
		Sort:     "email",
		Order:    pagination.SortAsc,
	})
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, int64(2), totalItems)
	assert.Equal(t, "anna@example.com", users[0].Email)
}

func TestUserAdapter_ListUsersCursor_AfterWithDuplicateSortValues(t *testing.T) {
	db, adapter := setupAdminUsersTestDB(t)
	sameCreatedAt := time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC)
	seedAdminUsers(t, db, []*models.User{
		{Model: gorm.Model{CreatedAt: sameCreatedAt}, Username: "anna", Email: "anna@example.com", DisplayName: "Anna Doe", PasswordHash: "hash"},
		{Model: gorm.Model{CreatedAt: sameCreatedAt}, Username: "bella", Email: "bella@example.com", DisplayName: "Bella Doe", PasswordHash: "hash"},
		{Model: gorm.Model{CreatedAt: sameCreatedAt.Add(time.Minute)}, Username: "carla", Email: "carla@example.com", DisplayName: "Carla Doe", PasswordHash: "hash"},
	})

	firstPage, err := adapter.ListUsersCursor(pagination.CursorQuery{
		PageSize: 2,
		Sort:     "created_at",
		Order:    pagination.SortAsc,
	})
	require.NoError(t, err)
	require.Len(t, firstPage.Users, 2)
	require.NotNil(t, firstPage.NextCursor)
	assert.Equal(t, "anna", firstPage.Users[0].Username)
	assert.Equal(t, "bella", firstPage.Users[1].Username)

	secondPage, err := adapter.ListUsersCursor(pagination.CursorQuery{
		PageSize: 2,
		Sort:     "created_at",
		Order:    pagination.SortAsc,
		After:    *firstPage.NextCursor,
	})
	require.NoError(t, err)
	require.Len(t, secondPage.Users, 1)
	assert.Equal(t, "carla", secondPage.Users[0].Username)
	assert.True(t, secondPage.HasPrev)
}

func TestUserAdapter_ListUsersCursor_BeforePreservesVisualOrder(t *testing.T) {
	db, adapter := setupAdminUsersTestDB(t)
	seedAdminUsers(t, db, []*models.User{
		{Username: "anna", Email: "anna@example.com", DisplayName: "Anna Doe", PasswordHash: "hash"},
		{Username: "bella", Email: "bella@example.com", DisplayName: "Bella Doe", PasswordHash: "hash"},
		{Username: "carla", Email: "carla@example.com", DisplayName: "Carla Doe", PasswordHash: "hash"},
	})

	pageOne, err := adapter.ListUsersCursor(pagination.CursorQuery{
		PageSize: 2,
		Sort:     "email",
		Order:    pagination.SortAsc,
	})
	require.NoError(t, err)
	require.NotNil(t, pageOne.NextCursor)

	pageTwo, err := adapter.ListUsersCursor(pagination.CursorQuery{
		PageSize: 2,
		Sort:     "email",
		Order:    pagination.SortAsc,
		After:    *pageOne.NextCursor,
	})
	require.NoError(t, err)
	require.NotNil(t, pageTwo.PrevCursor)
	require.Len(t, pageTwo.Users, 1)

	previousPage, err := adapter.ListUsersCursor(pagination.CursorQuery{
		PageSize: 2,
		Sort:     "email",
		Order:    pagination.SortAsc,
		Before:   *pageTwo.PrevCursor,
	})
	require.NoError(t, err)
	require.Len(t, previousPage.Users, 2)
	assert.Equal(t, "anna", previousPage.Users[0].Username)
	assert.Equal(t, "bella", previousPage.Users[1].Username)
	assert.True(t, previousPage.HasNext)
}

func TestUserAdapter_ListUsersCursor_InvalidCursor(t *testing.T) {
	db, adapter := setupAdminUsersTestDB(t)
	seedAdminUsers(t, db, []*models.User{
		{Username: "anna", Email: "anna@example.com", DisplayName: "Anna Doe", PasswordHash: "hash"},
	})

	result, err := adapter.ListUsersCursor(pagination.CursorQuery{
		PageSize: 10,
		Sort:     "email",
		Order:    pagination.SortAsc,
		After:    "invalid@@@",
	})
	assert.Nil(t, result)
	assert.ErrorIs(t, err, pagination.ErrInvalidCursor)
}
