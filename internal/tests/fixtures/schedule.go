package fixtures

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type ScheduleFixture struct {
	ID        string
	RoomID    string
	StartTime string // "09:00"
	EndTime   string // "18:00"
}

func DefaultScheduleFixture() ScheduleFixture {
	return ScheduleFixture{
		ID:        uuid.New().String(),
		RoomID:    uuid.New().String(),
		StartTime: "09:00",
		EndTime:   "18:00",
	}
}

func CreateSchedule(t *testing.T, db *sql.DB, fixture ScheduleFixture) string {
	t.Helper()

	query := `
        INSERT INTO schedules (id, room_id, start_time, end_time)
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

func CreateScheduleDays(t *testing.T, db *sql.DB, scheduleID string, days []int) {
	t.Helper()

	query := `
        INSERT INTO schedule_days (schedule_id, day_of_week)
        VALUES ($1, $2)
    `

	for _, day := range days {
		_, err := db.Exec(query, scheduleID, day)
		require.NoError(t, err)
	}
}
