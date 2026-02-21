package service

import (
	"strconv"
	"testing"

	"gosveltekit/internal/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_UpdateProfile(t *testing.T) {
	authService, _, _, _, _, db := setupTest(t)
	user := createTestUser(t, db)
	userID := strconv.FormatUint(uint64(user.ID), 10)

	displayName := "Updated User"
	firstName := "Updated"
	lastName := "Name"

	profile, err := authService.UpdateProfile(userID, UpdateProfileInput{
		DisplayName: &displayName,
		FirstName:   &firstName,
		LastName:    &lastName,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated User", profile.DisplayName)
	assert.Equal(t, "Updated", profile.FirstName)
	assert.Equal(t, "Name", profile.LastName)

	freshProfile, err := authService.GetProfile(userID)
	require.NoError(t, err)
	assert.Equal(t, "Updated User", freshProfile.DisplayName)
	assert.Equal(t, "Updated", freshProfile.FirstName)
	assert.Equal(t, "Name", freshProfile.LastName)
}

func TestAuthService_ChangePassword_Success(t *testing.T) {
	authService, _, _, _, _, db := setupTest(t)
	user := createTestUser(t, db)
	userID := strconv.FormatUint(uint64(user.ID), 10)

	err := authService.ChangePassword(userID, ChangePasswordInput{
		CurrentPassword: "password123",
		NewPassword:     "Secur3!PassA",
		ConfirmPassword: "Secur3!PassA",
	})
	require.NoError(t, err)

	_, err = authService.Login("testuser", "password123", "127.0.0.1", "test-agent")
	assert.ErrorIs(t, err, ErrInvalidCredentials)

	resp, err := authService.Login("testuser", "Secur3!PassA", "127.0.0.1", "test-agent")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.SessionID)
}

func TestAuthService_ChangePassword_WrongCurrentPassword(t *testing.T) {
	authService, _, _, _, _, db := setupTest(t)
	user := createTestUser(t, db)
	userID := strconv.FormatUint(uint64(user.ID), 10)

	err := authService.ChangePassword(userID, ChangePasswordInput{
		CurrentPassword: "invalid-current",
		NewPassword:     "Secur3!PassA",
		ConfirmPassword: "Secur3!PassA",
	})

	assert.ErrorIs(t, err, ErrWrongPassword)
}

func TestAuthService_ChangePassword_ValidationError(t *testing.T) {
	authService, _, _, _, _, db := setupTest(t)
	user := createTestUser(t, db)
	userID := strconv.FormatUint(uint64(user.ID), 10)

	err := authService.ChangePassword(userID, ChangePasswordInput{
		CurrentPassword: "password123",
		NewPassword:     "short",
		ConfirmPassword: "short",
	})

	assert.ErrorIs(t, err, validation.ErrPasswordTooShort)
}

func TestAuthService_ListAndRevokeSessions(t *testing.T) {
	authService, _, _, _, _, db := setupTest(t)
	user := createTestUser(t, db)
	userID := strconv.FormatUint(uint64(user.ID), 10)

	firstSession, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent-a")
	require.NoError(t, err)

	secondSession, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent-b")
	require.NoError(t, err)

	sessions, err := authService.ListSessions(userID, secondSession.SessionID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sessions), 2)

	var hasCurrent bool
	for _, session := range sessions {
		if session.ID == secondSession.SessionID {
			hasCurrent = session.IsCurrent
		}
	}
	assert.True(t, hasCurrent)

	err = authService.RevokeSession(userID, firstSession.SessionID, secondSession.SessionID)
	require.NoError(t, err)

	sessionsAfterRevoke, err := authService.ListSessions(userID, secondSession.SessionID)
	require.NoError(t, err)

	for _, session := range sessionsAfterRevoke {
		assert.NotEqual(t, firstSession.SessionID, session.ID)
	}
}
