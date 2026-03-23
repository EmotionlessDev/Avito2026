package slots

import (
	"context"
	"database/sql"
	"time"
)

type Slot struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type SlotStorage interface {
	CreateSlot(ctx context.Context, tx *sql.Tx, roomID string, startTime, endTime time.Time) (*Slot, error)
	GetFreeSlots(ctx context.Context, tx *sql.Tx, roomID string, dayStart, dayEnd time.Time) ([]*Slot, error)
	GetSlotsByDate(ctx context.Context, tx *sql.Tx, roomID string, dayStart, dayEnd time.Time) ([]*Slot, error)
}
