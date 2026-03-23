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
)

type GetSlots struct {
	scheduleStorage schedules.ScheduleStorage
	slotStorage     slots.SlotStorage
	roomStorage     rooms.RoomStorage

	db *sql.DB
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

	// tx
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

	// check if room exists
	_, err = uc.roomStorage.GetRoomByID(ctx, tx, input.RoomID)
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

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}

	if !contains(sched.DaysOfWeek, weekday) {
		return []*slots.Slot{}, nil
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	// check if slots already exist for this date
	existingSlots, err := uc.slotStorage.GetSlotsByDate(ctx, tx, input.RoomID, dayStart, dayEnd)
	if err != nil {
		return nil, fmt.Errorf("get slots: %w", err)
	}

	// if not, generate them based on schedule
	if len(existingSlots) == 0 {
		existingSlots, err = uc.generateSlots(ctx, tx, input.RoomID, sched, date)
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

func (uc *GetSlots) generateSlots(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	sched *schedules.Schedule,
	date time.Time,
) ([]*slots.Slot, error) {

	start, err := combineDateAndTime(date, sched.StartTime)
	if err != nil {
		return nil, err
	}

	end, err := combineDateAndTime(date, sched.EndTime)
	if err != nil {
		return nil, err
	}

	var result []*slots.Slot

	for t := start; t.Before(end); t = t.Add(30 * time.Minute) {
		slotEnd := t.Add(30 * time.Minute)

		slot, err := uc.slotStorage.CreateSlot(ctx, tx, roomID, t, slotEnd)
		if err != nil {
			return nil, fmt.Errorf("create slot: %w", err)
		}

		result = append(result, slot)
	}

	return result, nil
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func combineDateAndTime(date time.Time, timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}
	// set date to UTC to avoid timezone issues
	date = date.UTC()

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		time.UTC,
	), nil
}
