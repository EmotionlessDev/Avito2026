package slots

import (
	"context"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
)

type Slot struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type SlotStorage interface {
	CreateSlot(ctx context.Context, roomID string, startTime, endTime time.Time) (*Slot, error)
	GetFreeSlots(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]*Slot, error)
	GetSlotsByDate(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]*Slot, error)
	GetSlotByID(ctx context.Context, slotID string) (*Slot, error)
	CreateSlotsForSchedule(ctx context.Context, roomID string, sched *schedules.Schedule, startDate time.Time) ([]*Slot, error)
}
