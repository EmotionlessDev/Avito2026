package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
)

type GetRooms struct {
	roomStorage rooms.RoomStorage

	db *sql.DB
}

func NewGetRooms(roomStorage rooms.RoomStorage, db *sql.DB) *GetRooms {
	return &GetRooms{
		roomStorage: roomStorage,
		db:          db,
	}
}

func (uc *GetRooms) Execute(ctx context.Context) ([]*rooms.Room, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}

	tx, err := uc.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, common.ErrBeginTx
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	roomList, err := uc.roomStorage.GetRooms(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, common.ErrCommitTx
	}

	commited = true
	return roomList, nil
}
