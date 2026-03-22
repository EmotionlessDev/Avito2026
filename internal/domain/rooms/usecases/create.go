package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
)

type CreateRoom struct {
	roomStorage rooms.RoomStorage

	db *sql.DB
}

func NewCreateRoom(roomStorage rooms.RoomStorage, db *sql.DB) *CreateRoom {
	return &CreateRoom{
		roomStorage: roomStorage,
		db:          db,
	}
}

func (uc *CreateRoom) Execute(ctx context.Context, name, description string, capacity int) (string, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}

	tx, err := uc.db.BeginTx(ctx, opts)
	if err != nil {
		return "", common.ErrBeginTx
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	id, err := uc.roomStorage.CreateRoom(ctx, tx, name, description, capacity)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateRoom) {
			return "", common.ErrDuplicateRoom
		}
		return "", fmt.Errorf("failed to create room: %w", err)
	}

	return id, nil
}
