package rooms

import (
	"context"
	"database/sql"
)

type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	CreatedAt   string `json:"created_at"`
}

type RoomStorage interface {
	CreateRoom(ctx context.Context, tx *sql.Tx, name, description string, capacity int) (string, error)
}
