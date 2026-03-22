package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
	"github.com/lib/pq"
)

type pgRoom struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Capacity    int    `db:"capacity"`
	CreatedAt   string `db:"created_at"`
}

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

const createSQL = `
INSERT INTO rooms (name, description, capacity)
VALUES ($1, $2, $3)
RETURNING id
`

func (s *Storage) CreateRoom(ctx context.Context, tx *sql.Tx, name, description string, capacity int) (string, error) {
	if tx == nil {
		return "", common.ErrNilTx
	}

	var id string
	err := tx.QueryRowContext(ctx, createSQL, name, description, capacity).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "uniq_violation" {
			return "", common.ErrDuplicateRoom
		}
		return "", fmt.Errorf("failed to create room: %w", err)
	}

	return id, nil
}

const getRoomsSQL = `
	SELECT id, name, description, capacity, created_at
	FROM rooms
`

func (s *Storage) GetRooms(ctx context.Context, tx *sql.Tx) ([]*rooms.Room, error) {
	if tx == nil {
		return nil, common.ErrNilTx
	}

	rows, err := tx.QueryContext(ctx, getRoomsSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query rooms: %w", err)
	}
	defer rows.Close()

	var result []*rooms.Room
	for rows.Next() {
		var r pgRoom
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.Capacity, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		result = append(result, pgRoomToDomain(&r))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func pgRoomToDomain(r *pgRoom) *rooms.Room {
	return &rooms.Room{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Capacity:    r.Capacity,
		CreatedAt:   r.CreatedAt,
	}
}
