package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
)

type pgSlot struct {
	ID        string    `db:"id"`
	RoomID    string    `db:"room_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

const createSlotSQL = `
INSERT INTO slots (room_id, start_time, end_time)
VALUES ($1, $2, $3)
ON CONFLICT (room_id, start_time, end_time) DO NOTHING
RETURNING id, room_id, start_time, end_time
`

func (s *Storage) CreateSlot(
	ctx context.Context,
	roomID string,
	start, end time.Time,
) (*slots.Slot, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	start = start.UTC()
	end = end.UTC()

	var ps pgSlot

	err = tx.QueryRowContext(ctx, createSlotSQL, roomID, start, end).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return s.getSlotByTimeWithTx(ctx, tx, roomID, start, end)
		}
		return nil, fmt.Errorf("create slot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true
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
	roomID string,
	dayStart, dayEnd time.Time,
) ([]*slots.Slot, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	dayStart = dayStart.UTC()
	dayEnd = dayEnd.UTC()

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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true

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
AND s.start_time >= NOW()
AND b.id IS NULL
ORDER BY s.start_time
`

func (s *Storage) GetFreeSlots(
	ctx context.Context,
	roomID string,
	dayStart, dayEnd time.Time,
) ([]*slots.Slot, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	dayStart = dayStart.UTC()
	dayEnd = dayEnd.UTC()

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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true

	return result, nil
}

const getSlotByTimeSQL = `
SELECT id, room_id, start_time, end_time
FROM slots
WHERE room_id = $1 AND start_time = $2 AND end_time = $3
`

func (s *Storage) getSlotByTime(
	ctx context.Context,
	roomID string,
	start, end time.Time,
) (*slots.Slot, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	var ps pgSlot
	start = start.UTC()
	end = end.UTC()

	err = tx.QueryRowContext(ctx, getSlotByTimeSQL, roomID, start, end).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)
	if err != nil {
		return nil, fmt.Errorf("get slot by time: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true

	return pgSlotToDomain(&ps), nil
}

const getSlotByIDSQL = `
SELECT id, room_id, start_time, end_time
FROM slots
WHERE id = $1
`

func (s *Storage) GetSlotByID(
	ctx context.Context,
	slotID string,
) (*slots.Slot, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	var ps pgSlot

	err = tx.QueryRowContext(ctx, getSlotByIDSQL, slotID).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)
	if err != nil {
		return nil, fmt.Errorf("get slot by id: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true

	return pgSlotToDomain(&ps), nil
}

func (s *Storage) CreateSlotsForSchedule(ctx context.Context, roomID string, sched *schedules.Schedule, date time.Time) ([]*slots.Slot, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()
	var result []*slots.Slot
	date = date.UTC().Truncate(24 * time.Hour)

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	if !contains(sched.DaysOfWeek, weekday) {
		return []*slots.Slot{}, nil
	}

	startTime, err := combineDateAndTime(date, sched.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}
	endTime, err := combineDateAndTime(date, sched.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %w", err)
	}

	for t := startTime; t.Before(endTime); t = t.Add(30 * time.Minute) {
		slot, err := s.createSlotWithTx(ctx, tx, roomID, t, t.Add(30*time.Minute))
		if err != nil {
			return nil, err
		}
		result = append(result, slot)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true

	return result, nil
}

func (s *Storage) createSlotWithTx(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	start, end time.Time,
) (*slots.Slot, error) {
	start = start.UTC()
	end = end.UTC()

	var ps pgSlot
	err := tx.QueryRowContext(ctx, createSlotSQL, roomID, start, end).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return s.getSlotByTimeWithTx(ctx, tx, roomID, start, end)
		}
		return nil, fmt.Errorf("create slot: %w", err)
	}
	return pgSlotToDomain(&ps), nil
}

func (s *Storage) getSlotByTimeWithTx(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	start, end time.Time,
) (*slots.Slot, error) {
	var ps pgSlot
	start = start.UTC()
	end = end.UTC()

	err := tx.QueryRowContext(ctx, getSlotByTimeSQL, roomID, start, end).
		Scan(&ps.ID, &ps.RoomID, &ps.StartTime, &ps.EndTime)
	if err != nil {
		return nil, fmt.Errorf("get slot by time: %w", err)
	}
	return pgSlotToDomain(&ps), nil
}

func combineDateAndTime(date time.Time, timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}
	// set date to UTC to avoid timezone issues
	date = date.UTC()

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		time.UTC,
	), nil
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func pgSlotToDomain(s *pgSlot) *slots.Slot {
	return &slots.Slot{
		ID:        s.ID,
		RoomID:    s.RoomID,
		StartTime: s.StartTime.UTC().Format(time.RFC3339),
		EndTime:   s.EndTime.UTC().Format(time.RFC3339),
	}
}
