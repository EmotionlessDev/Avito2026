package fixtures

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type UserFixture struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
}

func DefaultUserFixture() UserFixture {
	id := uuid.New().String()
	return UserFixture{
		ID:           id,
		Email:        "user-" + id[:8] + "@example.com",
		PasswordHash: "dummyhash",
		Role:         "user",
	}
}

func AdminUserFixture() UserFixture {
	id := uuid.New().String()
	return UserFixture{
		ID:           id,
		Email:        "admin-" + id[:8] + "@example.com",
		PasswordHash: "dummyhash",
		Role:         "admin",
	}
}

func CreateUser(t *testing.T, db *sql.DB, fixture UserFixture) string {
	t.Helper()

	query := `
        INSERT INTO users (id, email, password_hash, role)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id string
	err := db.QueryRow(
		query,
		fixture.ID,
		fixture.Email,
		fixture.PasswordHash,
		fixture.Role,
	).Scan(&id)
	require.NoError(t, err)

	return id
}

func GetTestUser(t *testing.T, db *sql.DB, role string) string {
	t.Helper()

	query := `SELECT id FROM users WHERE role = $1 LIMIT 1`

	var id string
	err := db.QueryRow(query, role).Scan(&id)
	require.NoError(t, err)

	return id
}
