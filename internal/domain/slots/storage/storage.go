package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
)

type pgSlot struct {
	ID        string    `db:"id"`
	RoomID    string    `db:"room_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

const createSlotSQL = `
INSERT INTO slots (room_id, start_time, end_time)
VALUES ($1, $2, $3)
ON CONFLICT (room_id, start_time, end_time) DO NOTHING
RETURNING id, room_id, start_time, end_time
`

func (s *Storage) CreateSlot(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	start, end time.Time,
) (*slots.Slot, error) {

	if tx == nil {
		return nil, common.ErrNilTx
	}

	var ps pgSlot

	err := tx.QueryRowContext(ctx, createSlotSQL, roomID, start.UTC(), end.UTC()).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return s.getSlotByTime(ctx, tx, roomID, start, end)
		}
		return nil, fmt.Errorf("create slot: %w", err)
	}

	return pgSlotToDomain(&ps), nil
}

const getSlotsByDateSQL = `
SELECT id, room_id, start_time, end_time
FROM slots
WHERE room_id = $1
AND start_time >= $2
AND start_time < $3
ORDER BY start_time
`

func (s *Storage) GetSlotsByDate(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	dayStart, dayEnd time.Time,
) ([]*slots.Slot, error) {

	if tx == nil {
		return nil, common.ErrNilTx
	}

	rows, err := tx.QueryContext(ctx, getSlotsByDateSQL, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, fmt.Errorf("get slots: %w", err)
	}
	defer rows.Close()

	var result []*slots.Slot

	for rows.Next() {
		var ps pgSlot

		if err := rows.Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime); err != nil {
			return nil, fmt.Errorf("scan slot: %w", err)
		}

		result = append(result, pgSlotToDomain(&ps))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

const getFreeSlotsSQL = `
SELECT s.id, s.room_id, s.start_time, s.end_time
FROM slots s
LEFT JOIN bookings b
  ON b.slot_id = s.id AND b.status = 'active'
WHERE s.room_id = $1
AND s.start_time >= $2
AND s.start_time < $3
AND b.id IS NULL
ORDER BY s.start_time
`

func (s *Storage) GetFreeSlots(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	dayStart, dayEnd time.Time,
) ([]*slots.Slot, error) {

	if tx == nil {
		return nil, common.ErrNilTx
	}

	rows, err := tx.QueryContext(ctx, getFreeSlotsSQL, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, fmt.Errorf("get free slots: %w", err)
	}
	defer rows.Close()

	var result []*slots.Slot

	for rows.Next() {
		var ps pgSlot

		if err := rows.Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime); err != nil {
			return nil, fmt.Errorf("scan free slot: %w", err)
		}

		result = append(result, pgSlotToDomain(&ps))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

const getSlotByTimeSQL = `
SELECT id, room_id, start_time, end_time
FROM slots
WHERE room_id = $1 AND start_time = $2 AND end_time = $3
`

func (s *Storage) getSlotByTime(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	start, end time.Time,
) (*slots.Slot, error) {

	var ps pgSlot

	err := tx.QueryRowContext(ctx, getSlotByTimeSQL, roomID, start, end).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)
	if err != nil {
		return nil, fmt.Errorf("get slot by time: %w", err)
	}

	return pgSlotToDomain(&ps), nil
}

func pgSlotToDomain(s *pgSlot) *slots.Slot {
	return &slots.Slot{
		ID:        s.ID,
		RoomID:    s.RoomID,
		StartTime: s.StartTime.UTC().Format(time.RFC3339),
		EndTime:   s.EndTime.UTC().Format(time.RFC3339),
	}
}
