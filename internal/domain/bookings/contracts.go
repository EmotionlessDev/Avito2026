package bookings

import (
	"context"
	"database/sql"
)

type Booking struct {
	ID             string `json:"id"`
	SlotID         string `json:"slot_id"`
	UserID         string `json:"user_id"`
	Status         string `json:"status"`
	ConferenceLink string `json:"conference_link,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type BookingStorage interface {
	CreateBooking(ctx context.Context, tx *sql.Tx, slotID, userID string, conferenceLink *string) (*Booking, error)
	GetBookingsPaginated(ctx context.Context, tx *sql.Tx, limit, offset int) ([]*Booking, int, error)
	GetBookingsByUserID(ctx context.Context, tx *sql.Tx, userID string) ([]*Booking, error)
	GetBookingByID(ctx context.Context, tx *sql.Tx, bookingID string) (*Booking, error)
	UpdateBookingStatus(ctx context.Context, tx *sql.Tx, bookingID, status string) error
}
