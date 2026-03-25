package fixtures

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type SlotFixture struct {
	ID        string
	RoomID    string
	StartTime time.Time
	EndTime   time.Time
}

func DefaultSlotFixture() SlotFixture {
	now := time.Now().UTC()
	return SlotFixture{
		ID:        uuid.New().String(),
		RoomID:    uuid.New().String(),
		StartTime: now.Add(24 * time.Hour),
		EndTime:   now.Add(24*time.Hour + 30*time.Minute),
	}
}

func FutureSlotFixture(roomID string, date time.Time) SlotFixture {
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
	return SlotFixture{
		ID:        uuid.New().String(),
		RoomID:    roomID,
		StartTime: startTime,
		EndTime:   startTime.Add(30 * time.Minute),
	}
}

func CreateSlot(t *testing.T, db *sql.DB, fixture SlotFixture) string {
	t.Helper()

	query := `
        INSERT INTO slots (id, room_id, start_time, end_time)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id string
	err := db.QueryRow(
		query,
		fixture.ID,
		fixture.RoomID,
		fixture.StartTime,
		fixture.EndTime,
	).Scan(&id)
	require.NoError(t, err)

	return id
}

