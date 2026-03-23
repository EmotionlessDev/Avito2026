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
	if tx == nil {
		return nil, common.ErrNilTx
	}

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

const getBookingsSQL = `
SELECT id, slot_id, user_id, status, conference_link, created_at
FROM bookings
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

const countBookingsSQL = `
SELECT COUNT(*) FROM bookings
`

func (s *Storage) GetBookingsPaginated(
	ctx context.Context,
	tx *sql.Tx,
	limit, offset int,
) ([]*bookings.Booking, int, error) {

	if tx == nil {
		return nil, 0, common.ErrNilTx
	}

	// total count
	var total int
	err := tx.QueryRowContext(ctx, countBookingsSQL).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	rows, err := tx.QueryContext(ctx, getBookingsSQL, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer rows.Close()

	var result []*bookings.Booking

	for rows.Next() {
		var pb pgBooking

		err := rows.Scan(
			&pb.ID,
			&pb.SlotID,
			&pb.UserID,
			&pb.Status,
			&pb.ConferenceLink,
			&pb.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan booking: %w", err)
		}

		result = append(result, pgBookingToDomain(&pb))
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

const getBookingsByUserIDSQL = `
SELECT 
    b.id, 
    b.slot_id, 
    b.user_id, 
    b.status, 
    b.conference_link, 
    b.created_at
FROM bookings b
JOIN slots s ON b.slot_id = s.id
WHERE 
    b.user_id = $1 
    AND s.start_time >= NOW()
    AND b.status = 'active'
ORDER BY s.start_time ASC
`

func (s *Storage) GetBookingsByUserID(
	ctx context.Context,
	tx *sql.Tx,
	userID string,
) ([]*bookings.Booking, error) {
	if tx == nil {
		return nil, common.ErrNilTx
	}

	rows, err := tx.QueryContext(ctx, getBookingsByUserIDSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by user id: %w", err)
	}
	defer rows.Close()

	var result []*bookings.Booking
	for rows.Next() {
		var pb pgBooking

		err := rows.Scan(
			&pb.ID,
			&pb.SlotID,
			&pb.UserID,
			&pb.Status,
			&pb.ConferenceLink,
			&pb.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}

		result = append(result, pgBookingToDomain(&pb))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
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
