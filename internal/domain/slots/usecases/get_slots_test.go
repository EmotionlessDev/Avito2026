package usecases

import (
	"context"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// type GetSlots struct {
// 	scheduleStorage schedules.ScheduleStorage
// 	slotStorage     slots.SlotStorage
// 	roomStorage     rooms.RoomStorage
// }
//
// func NewGetSlots(
// 	scheduleStorage schedules.ScheduleStorage,
// 	slotStorage slots.SlotStorage,
// 	roomStorage rooms.RoomStorage,
// ) *GetSlots {
// 	return &GetSlots{
// 		scheduleStorage: scheduleStorage,
// 		slotStorage:     slotStorage,
// 		roomStorage:     roomStorage,
// 	}
// }
//
// type GetSlotsInput struct {
// 	RoomID string
// 	Date   string // "YYYY-MM-DD"
// }
//
// func (uc *GetSlots) Execute(ctx context.Context, input GetSlotsInput) ([]*slots.Slot, error) {
// 	// parse date
// 	date, err := time.Parse("2006-01-02", input.Date)
// 	if err != nil {
// 		return nil, common.ErrInvalidDate
// 	}
// 	date = date.UTC()
//
// 	// check if date is in the past
// 	if date.Before(time.Now().UTC().Truncate(24 * time.Hour)) {
// 		return []*slots.Slot{}, nil
// 	}
//
// 	// check uuid
// 	if _, err := uuid.Parse(input.RoomID); err != nil {
// 		return nil, common.ErrInvalidUUID
// 	}
//
// 	// check room exists
// 	_, err = uc.roomStorage.GetRoomByID(ctx, input.RoomID)
// 	if err != nil {
// 		if err == common.ErrRoomNotFound {
// 			return nil, common.ErrRoomNotFound
// 		}
// 		return nil, fmt.Errorf("get room: %w", err)
// 	}
//
// 	// get schedule
// 	sched, err := uc.scheduleStorage.GetScheduleByRoomID(ctx, input.RoomID)
// 	if err != nil {
// 		if err == common.ErrScheduleNotFound {
// 			return []*slots.Slot{}, nil
// 		}
// 		return nil, fmt.Errorf("get schedule: %w", err)
// 	}
//
// 	// check day of week
// 	weekday := int(date.Weekday())
// 	if weekday == 0 {
// 		weekday = 7
// 	}
// 	if !contains(sched.DaysOfWeek, weekday) {
// 		return []*slots.Slot{}, nil
// 	}
//
// 	dayStart := date
// 	dayEnd := dayStart.Add(24 * time.Hour)
//
// 	// get existing slots
// 	existingSlots, err := uc.slotStorage.GetSlotsByDate(ctx, input.RoomID, dayStart, dayEnd)
// 	if err != nil {
// 		return nil, fmt.Errorf("get slots: %w", err)
// 	}
//
// 	// generate slots if none exist
// 	if len(existingSlots) == 0 {
// 		existingSlots, err = uc.slotStorage.CreateSlotsForSchedule(ctx, input.RoomID, sched, date)
// 		if err != nil {
// 			return nil, fmt.Errorf("generate slots: %w", err)
// 		}
// 	}
//
// 	// get free slots
// 	freeSlots, err := uc.slotStorage.GetFreeSlots(ctx, input.RoomID, dayStart, dayEnd)
// 	if err != nil {
// 		return nil, fmt.Errorf("get free slots: %w", err)
// 	}
//
// 	return freeSlots, nil
// }
//
// func contains(slice []int, item int) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }

const (
	// Dates
	testFutureDate = "2220-06-10"
	testPastDate   = "2020-01-01"

	// UUID
	testRoomID     = "550e8400-e29b-41d4-a716-446655440000"
	testUserID     = "660e8400-e29b-41d4-a716-446655440001"
	testSlotID     = "770e8400-e29b-41d4-a716-446655440002"
	testScheduleID = "schedule-id"

	// Invalid uuid
	testInvalidUUID = "not-a-uuid"
)

func TestGetSlots_Execute_Success_ExistingSlots(t *testing.T) {
	scheduleStorageMock := mocks.NewMockScheduleStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	roomStorageMock := mocks.NewMockRoomStorage(t)

	uc := NewGetSlots(scheduleStorageMock, slotStorageMock, roomStorageMock)

	roomStorageMock.EXPECT().
		GetRoomByID(mock.Anything, testRoomID).
		Return(&rooms.Room{
			ID:   testRoomID,
			Name: "Test Room",
		}, nil)

	scheduleStorageMock.EXPECT().
		GetScheduleByRoomID(mock.Anything, testRoomID).
		Return(&schedules.Schedule{
			ID:         testScheduleID,
			RoomID:     testRoomID,
			DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7},
			StartTime:  "09:00",
			EndTime:    "10:00",
		}, nil)

	slotStorageMock.EXPECT().
		GetSlotsByDate(mock.Anything, testRoomID, mock.Anything, mock.Anything).
		Return([]*slots.Slot{
			{ID: "slot-id-1", RoomID: testRoomID, StartTime: "2220-06-10T09:00:00Z", EndTime: "2220-06-10T09:30:00Z"},
			{ID: "slot-id-2", RoomID: testRoomID, StartTime: "2220-06-10T09:30:00Z", EndTime: "2220-06-10T10:00:00Z"},
		}, nil)

	slotStorageMock.EXPECT().
		GetFreeSlots(mock.Anything, testRoomID, mock.Anything, mock.Anything).
		Return([]*slots.Slot{
			{ID: "free-slot-1", RoomID: testRoomID, StartTime: "2220-06-10T09:00:00Z", EndTime: "2220-06-10T09:30:00Z"},
		}, nil)

	input := GetSlotsInput{
		RoomID: testRoomID,
		Date:   testFutureDate,
	}

	slots, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Len(t, slots, 1)
	assert.Equal(t, "free-slot-1", slots[0].ID)
}

func TestGetSlots_Execute_Success_NoExistingSlots(t *testing.T) {
	scheduleStorageMock := mocks.NewMockScheduleStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	roomStorageMock := mocks.NewMockRoomStorage(t)

	uc := NewGetSlots(scheduleStorageMock, slotStorageMock, roomStorageMock)

	roomStorageMock.EXPECT().
		GetRoomByID(mock.Anything, testRoomID).
		Return(&rooms.Room{
			ID:   testRoomID,
			Name: "Test Room",
		}, nil)

	scheduleStorageMock.EXPECT().
		GetScheduleByRoomID(mock.Anything, testRoomID).
		Return(&schedules.Schedule{
			ID:         testScheduleID,
			RoomID:     testRoomID,
			DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7},
			StartTime:  "09:00",
			EndTime:    "10:00",
		}, nil)

	slotStorageMock.EXPECT().
		GetSlotsByDate(mock.Anything, testRoomID, mock.Anything, mock.Anything).
		Return([]*slots.Slot{}, nil)

	slotStorageMock.EXPECT().
		CreateSlotsForSchedule(mock.Anything, testRoomID, mock.Anything, mock.Anything).
		Return([]*slots.Slot{
			{ID: "slot-id-1", RoomID: testRoomID, StartTime: "2220-06-10T09:00:00Z", EndTime: "2220-06-10T09:30:00Z"},
			{ID: "slot-id-2", RoomID: testRoomID, StartTime: "2220-06-10T09:30:00Z", EndTime: "2220-06-10T10:00:00Z"},
		}, nil)

	slotStorageMock.EXPECT().
		GetFreeSlots(mock.Anything, testRoomID, mock.Anything, mock.Anything).
		Return([]*slots.Slot{
			{ID: "free-slot-1", RoomID: testRoomID, StartTime: "2220-06-10T09:00:00Z", EndTime: "2220-06-10T09:30:00Z"},
			{ID: "free-slot-2", RoomID: testRoomID, StartTime: "2220-06-10T09:30:00Z", EndTime: "2220-06-10T10:00:00Z"},
		}, nil)

	input := GetSlotsInput{
		RoomID: testRoomID,
		Date:   testFutureDate,
	}

	slots, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Len(t, slots, 2)
	assert.Equal(t, "free-slot-1", slots[0].ID)
	assert.Equal(t, "free-slot-2", slots[1].ID)
}

func TestGetSlots_Execute_Error_InvalidDate(t *testing.T) {
	scheduleStorageMock := mocks.NewMockScheduleStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	roomStorageMock := mocks.NewMockRoomStorage(t)

	uc := NewGetSlots(scheduleStorageMock, slotStorageMock, roomStorageMock)

	input := GetSlotsInput{
		RoomID: testRoomID,
		Date:   "invalid-date",
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, common.ErrInvalidDate, err)
}

func TestGetSlots_Execute_Error_InvalidUUID(t *testing.T) {
	scheduleStorageMock := mocks.NewMockScheduleStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	roomStorageMock := mocks.NewMockRoomStorage(t)

	uc := NewGetSlots(scheduleStorageMock, slotStorageMock, roomStorageMock)

	input := GetSlotsInput{
		RoomID: testInvalidUUID,
		Date:   testFutureDate,
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, common.ErrInvalidUUID, err)
}

func TestGetSlots_Execute_WeekdayNotInSchedule(t *testing.T) {
	scheduleStorageMock := mocks.NewMockScheduleStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	roomStorageMock := mocks.NewMockRoomStorage(t)

	uc := NewGetSlots(scheduleStorageMock, slotStorageMock, roomStorageMock)

	roomStorageMock.EXPECT().
		GetRoomByID(mock.Anything, testRoomID).
		Return(&rooms.Room{
			ID:   testRoomID,
			Name: "Test Room",
		}, nil)

	scheduleStorageMock.EXPECT().
		GetScheduleByRoomID(mock.Anything, testRoomID).
		Return(&schedules.Schedule{
			ID:         testScheduleID,
			RoomID:     testRoomID,
			DaysOfWeek: []int{1},
			StartTime:  "09:00",
			EndTime:    "10:00",
		}, nil)

	input := GetSlotsInput{
		RoomID: testRoomID,
		Date:   "2220-06-13",
	}

	slots, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Len(t, slots, 0)
}

func TestGetSlots_Execute_Error_RoomNotFound(t *testing.T) {
	scheduleStorageMock := mocks.NewMockScheduleStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	roomStorageMock := mocks.NewMockRoomStorage(t)

	uc := NewGetSlots(scheduleStorageMock, slotStorageMock, roomStorageMock)

	roomStorageMock.EXPECT().
		GetRoomByID(mock.Anything, testRoomID).
		Return(nil, common.ErrRoomNotFound)

	input := GetSlotsInput{
		RoomID: testRoomID,
		Date:   testFutureDate,
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Equal(t, common.ErrRoomNotFound, err)
}
