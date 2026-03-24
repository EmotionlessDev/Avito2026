package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
	"github.com/google/uuid"
)

type GetSlots struct {
	scheduleStorage schedules.ScheduleStorage
	slotStorage     slots.SlotStorage
	roomStorage     rooms.RoomStorage
	db              *sql.DB
}

func NewGetSlots(
	scheduleStorage schedules.ScheduleStorage,
	slotStorage slots.SlotStorage,
	roomStorage rooms.RoomStorage,
	db *sql.DB,
) *GetSlots {
	return &GetSlots{
		scheduleStorage: scheduleStorage,
		slotStorage:     slotStorage,
		roomStorage:     roomStorage,
		db:              db,
	}
}

type GetSlotsInput struct {
	RoomID string
	Date   string // "YYYY-MM-DD"
}

func (uc *GetSlots) Execute(ctx context.Context, input GetSlotsInput) ([]*slots.Slot, error) {
	// parse date
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, common.ErrInvalidDate
	}
	date = date.UTC()

	// check if date is in the past
	if date.Before(time.Now().UTC().Truncate(24 * time.Hour)) {
		return []*slots.Slot{}, nil
	}

	// begin tx
	tx, err := uc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, common.ErrBeginTx
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// check uuid
	if _, err := uuid.Parse(input.RoomID); err != nil {
		return nil, common.ErrInvalidUUID
	}

	// check room exists
	_, err = uc.roomStorage.GetRoomByID(ctx, input.RoomID)
	if err != nil {
		if err == common.ErrRoomNotFound {
			return nil, common.ErrRoomNotFound
		}
		return nil, fmt.Errorf("get room: %w", err)
	}

	// get schedule
	sched, err := uc.scheduleStorage.GetScheduleByRoomID(ctx, tx, input.RoomID)
	if err != nil {
		if err == common.ErrScheduleNotFound {
			return []*slots.Slot{}, nil
		}
		return nil, fmt.Errorf("get schedule: %w", err)
	}

	// check day of week
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	if !contains(sched.DaysOfWeek, weekday) {
		return []*slots.Slot{}, nil
	}

	dayStart := date
	dayEnd := dayStart.Add(24 * time.Hour)

	// get existing slots
	existingSlots, err := uc.slotStorage.GetSlotsByDate(ctx, tx, input.RoomID, dayStart, dayEnd)
	if err != nil {
		return nil, fmt.Errorf("get slots: %w", err)
	}

	// generate slots if none exist
	if len(existingSlots) == 0 {
		existingSlots, err = uc.slotStorage.CreateSlotsForSchedule(ctx, tx, input.RoomID, sched, date)
		if err != nil {
			return nil, fmt.Errorf("generate slots: %w", err)
		}
	}

	// get free slots
	freeSlots, err := uc.slotStorage.GetFreeSlots(ctx, tx, input.RoomID, dayStart, dayEnd)
	if err != nil {
		return nil, fmt.Errorf("get free slots: %w", err)
	}

	// commit
	if err := tx.Commit(); err != nil {
		return nil, common.ErrCommitTx
	}
	committed = true

	return freeSlots, nil
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
