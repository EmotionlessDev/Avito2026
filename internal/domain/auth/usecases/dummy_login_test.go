package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
)

const (
	testJWTSecret = "test-secret-key-123"
	testAdminID   = "11111111-1111-1111-1111-111111111111"
	testUserID    = "22222222-2222-2222-2222-222222222222"
)

func TestDummyLogin_Execute_Success_Admin(t *testing.T) {
	uc := NewDummyLogin(testJWTSecret)

	response, err := uc.Execute(context.Background(), "admin")

	require.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, testAdminID, response.UserID)
	assert.Equal(t, "admin", response.Role)
	assert.WithinDuration(t, time.Now().UTC(), response.CreatedAt, 5*time.Second)
}

func TestDummyLogin_Execute_Success_User(t *testing.T) {
	uc := NewDummyLogin(testJWTSecret)

	response, err := uc.Execute(context.Background(), "user")

	require.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, testUserID, response.UserID)
	assert.Equal(t, "user", response.Role)
	assert.WithinDuration(t, time.Now().UTC(), response.CreatedAt, 5*time.Second)
}

func TestDummyLogin_Execute_Error_InvalidRole(t *testing.T) {
	uc := NewDummyLogin(testJWTSecret)

	response, err := uc.Execute(context.Background(), "invalid-role")

	require.Error(t, err)
	assert.ErrorIs(t, err, common.ErrInvalidRole)
	assert.Empty(t, response.Token)
	assert.Empty(t, response.UserID)
	assert.Empty(t, response.Role)
}

func TestDummyLogin_Execute_Error_EmptyRole(t *testing.T) {
	uc := NewDummyLogin(testJWTSecret)

	response, err := uc.Execute(context.Background(), "")

	require.Error(t, err)
	assert.ErrorIs(t, err, common.ErrInvalidRole)
	assert.Empty(t, response.Token)
	assert.Empty(t, response.UserID)
}

func TestDummyLogin_Execute_TokenFormat(t *testing.T) {
	uc := NewDummyLogin(testJWTSecret)

	response, err := uc.Execute(context.Background(), "user")

	require.NoError(t, err)

	assert.Contains(t, response.Token, ".")

	parts := splitToken(response.Token)
	require.Len(t, parts, 3)
	assert.NotEmpty(t, parts[0]) // header
	assert.NotEmpty(t, parts[1]) // payload
	assert.NotEmpty(t, parts[2]) // signature
}

func TestDummyLogin_Execute_DifferentSecrets(t *testing.T) {
	uc1 := NewDummyLogin("secret-1")
	uc2 := NewDummyLogin("secret-2")

	response1, err1 := uc1.Execute(context.Background(), "user")
	response2, err2 := uc2.Execute(context.Background(), "user")

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.NotEqual(t, response1.Token, response2.Token)

	assert.Equal(t, response1.UserID, response2.UserID)
}

func TestDummyLogin_Execute_ConsistentUserID(t *testing.T) {
	uc := NewDummyLogin(testJWTSecret)

	response1, err1 := uc.Execute(context.Background(), "admin")
	response2, err2 := uc.Execute(context.Background(), "admin")

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.Equal(t, response1.UserID, response2.UserID)
	assert.Equal(t, testAdminID, response1.UserID)
	assert.Equal(t, testAdminID, response2.UserID)
}

func TestDummyLogin_Execute_RoleMapping(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		wantUserID string
		wantErr    error
	}{
		{
			name:       "admin role",
			role:       "admin",
			wantUserID: testAdminID,
			wantErr:    nil,
		},
		{
			name:       "user role",
			role:       "user",
			wantUserID: testUserID,
			wantErr:    nil,
		},
		{
			name:       "invalid role",
			role:       "manager",
			wantUserID: "",
			wantErr:    common.ErrInvalidRole,
		},
		{
			name:       "empty role",
			role:       "",
			wantUserID: "",
			wantErr:    common.ErrInvalidRole,
		},
		{
			name:       "case sensitive",
			role:       "Admin",
			wantUserID: "",
			wantErr:    common.ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewDummyLogin(testJWTSecret)

			response, err := uc.Execute(context.Background(), tt.role)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, response.UserID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserID, response.UserID)
				assert.Equal(t, tt.role, response.Role)
				assert.NotEmpty(t, response.Token)
			}
		})
	}
}

func splitToken(token string) []string {
	parts := make([]string, 0)
	start := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	parts = append(parts, token[start:])
	return parts
}
