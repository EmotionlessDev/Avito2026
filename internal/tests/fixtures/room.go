package fixtures

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type RoomFixture struct {
	ID          string
	Name        string
	Description string
	Capacity    int
}

func DefaultRoomFixture() RoomFixture {
	return RoomFixture{
		ID:          uuid.New().String(),
		Name:        "Test Room " + uuid.New().String()[:8],
		Description: "Test description",
		Capacity:    10,
	}
}

func CreateRoom(t *testing.T, db *sql.DB, fixture RoomFixture) string {
	t.Helper()

	query := `
        INSERT INTO rooms (id, name, description, capacity)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id string
	err := db.QueryRow(
		query,
		fixture.ID,
		fixture.Name,
		fixture.Description,
		fixture.Capacity,
	).Scan(&id)
	require.NoError(t, err)

	return id
}
