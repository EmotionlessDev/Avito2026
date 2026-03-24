package bookings

import (
	"context"
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
	CreateBooking(ctx context.Context, slotID, userID string, conferenceLink *string) (*Booking, error)
	GetBookingsPaginated(ctx context.Context, limit, offset int) ([]*Booking, int, error)
	GetBookingsByUserID(ctx context.Context, userID string) ([]*Booking, error)
	GetBookingByID(ctx context.Context, bookingID string) (*Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID, status string) error
}
