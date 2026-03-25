package fixtures

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type BookingFixture struct {
	ID             string
	SlotID         string
	UserID         string
	Status         string
	ConferenceLink string
}

func DefaultBookingFixture() BookingFixture {
	return BookingFixture{
		ID:             uuid.New().String(),
		SlotID:         uuid.New().String(),
		UserID:         uuid.New().String(),
		Status:         "active",
		ConferenceLink: "https://meet.example.com/room-" + uuid.New().String()[:8],
	}
}

func CreateBooking(t *testing.T, db *sql.DB, fixture BookingFixture) string {
	t.Helper()

	query := `
        INSERT INTO bookings (id, slot_id, user_id, status, conference_link)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	var id string
	err := db.QueryRow(
		query,
		fixture.ID,
		fixture.SlotID,
		fixture.UserID,
		fixture.Status,
		fixture.ConferenceLink,
	).Scan(&id)
	require.NoError(t, err)

	return id
}

func CreateBookingWithTx(t *testing.T, tx *sql.Tx, fixture BookingFixture) string {
	t.Helper()

	query := `
        INSERT INTO bookings (id, slot_id, user_id, status, conference_link)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	var id string
	err := tx.QueryRow(
		query,
		fixture.ID,
		fixture.SlotID,
		fixture.UserID,
		fixture.Status,
		fixture.ConferenceLink,
	).Scan(&id)
	require.NoError(t, err)

	return id
}

