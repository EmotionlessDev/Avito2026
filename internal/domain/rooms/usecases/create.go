package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
)

type CreateRoom struct {
	roomStorage rooms.RoomStorage
}

func NewCreateRoom(roomStorage rooms.RoomStorage) *CreateRoom {
	return &CreateRoom{
		roomStorage: roomStorage,
	}
}

func (uc *CreateRoom) Execute(ctx context.Context, name, description string, capacity int) (string, error) {
	id, err := uc.roomStorage.CreateRoom(ctx, name, description, capacity)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateRoom) {
			return "", common.ErrDuplicateRoom
		}
		return "", fmt.Errorf("failed to create room: %w", err)
	}

	return id, nil
}
