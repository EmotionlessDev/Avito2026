package rooms

import (
	"context"
)

type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	CreatedAt   string `json:"created_at"`
}

type RoomStorage interface {
	CreateRoom(ctx context.Context, name, description string, capacity int) (string, error)
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	GetRooms(ctx context.Context) ([]*Room, error)
}
