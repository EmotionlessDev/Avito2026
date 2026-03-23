package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/lib/pq"
)

type pgBooking struct {
	ID             string         `db:"id"`
	SlotID         string         `db:"slot_id"`
	UserID         string         `db:"user_id"`
	Status         string         `db:"status"`
	ConferenceLink sql.NullString `db:"conference_link"`
	CreatedAt      string         `db:"created_at"`
}

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

const createBookingSQL = `
INSERT INTO bookings (slot_id, user_id, conference_link)
VALUES ($1, $2, $3)
RETURNING id, slot_id, user_id, status, conference_link, created_at
`

func (s *Storage) CreateBooking(
	ctx context.Context,
	tx *sql.Tx,
	slotID, userID string,
	conferenceLink *string,
) (*bookings.Booking, error) {

	var pb pgBooking

	err := tx.QueryRowContext(ctx, createBookingSQL, slotID, userID, conferenceLink).
		Scan(&pb.ID, &pb.SlotID, &pb.UserID, &pb.Status, &pb.ConferenceLink, &pb.CreatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return nil, common.ErrSlotAlreadyBooked
		}
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return pgBookingToDomain(&pb), nil
}

func pgBookingToDomain(pb *pgBooking) *bookings.Booking {
	var link string
	if pb.ConferenceLink.Valid {
		link = pb.ConferenceLink.String
	}

	return &bookings.Booking{
		ID:             pb.ID,
		SlotID:         pb.SlotID,
		UserID:         pb.UserID,
		Status:         pb.Status,
		ConferenceLink: link,
		CreatedAt:      pb.CreatedAt,
	}
}

func isUniqueViolation(err error) bool {
	// pq specific
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
