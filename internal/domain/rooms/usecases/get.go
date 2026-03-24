package usecases

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
)

type GetRooms struct {
	roomStorage rooms.RoomStorage
}

func NewGetRooms(roomStorage rooms.RoomStorage) *GetRooms {
	return &GetRooms{
		roomStorage: roomStorage,
	}
}

func (uc *GetRooms) Execute(ctx context.Context) ([]*rooms.Room, error) {
	roomList, err := uc.roomStorage.GetRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}

	return roomList, nil
}
