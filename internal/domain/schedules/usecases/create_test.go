package usecases

import (
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateSchedule_Execute_Success(t *testing.T) {
	storageMock := mocks.NewMockScheduleStorage(t)
	uc := NewCreateSchedule(storageMock)

	storageMock.EXPECT().IsScheduleExistsByRoomID(mock.Anything, "roomId").Return(false, nil)
	storageMock.EXPECT().CreateSchedule(mock.Anything, "roomId", mock.Anything, mock.Anything, []int{1, 3, 5}).Return(&schedules.Schedule{
		ID:         "scheduleId",
		RoomID:     "roomId",
		StartTime:  "09:00",
		EndTime:    "17:00",
		DaysOfWeek: []int{1, 3, 5},
	}, nil)

	input := dto.CreateScheduleInput{
		RoomID:     "roomId",
		StartTime:  "09:00",
		EndTime:    "17:00",
		DaysOfWeek: []int{1, 3, 5},
	}

	schedule, err := uc.Execute(nil, input)

	require.NoError(t, err)
	assert.Equal(t, "scheduleId", schedule.ID)
	assert.Equal(t, "roomId", schedule.RoomID)
	assert.Equal(t, "09:00", schedule.StartTime)
	assert.Equal(t, "17:00", schedule.EndTime)
	assert.Equal(t, []int{1, 3, 5}, schedule.DaysOfWeek)
}

func TestCreateSchedule_Execute_InvalidDaysOfWeek(t *testing.T) {
	storageMock := mocks.NewMockScheduleStorage(t)
	uc := NewCreateSchedule(storageMock)

	input := dto.CreateScheduleInput{
		RoomID:     "roomId",
		StartTime:  "09:00",
		EndTime:    "17:00",
		DaysOfWeek: []int{0, 8},
	}

	_, err := uc.Execute(nil, input)

	require.ErrorIs(t, err, common.ErrInvalidScheduleDay)
}

func TestCreateSchedule_Execute_DuplicateDaysOfWeek(t *testing.T) {
	storageMock := mocks.NewMockScheduleStorage(t)
	uc := NewCreateSchedule(storageMock)

	input := dto.CreateScheduleInput{
		RoomID:     "roomId",
		StartTime:  "09:00",
		EndTime:    "17:00",
		DaysOfWeek: []int{1, 1},
	}

	_, err := uc.Execute(nil, input)

	require.ErrorIs(t, err, common.ErrDuplicateScheduleDay)
}

func TestCreateSchedule_Execute_InvalidScheduleTime(t *testing.T) {
	storageMock := mocks.NewMockScheduleStorage(t)
	uc := NewCreateSchedule(storageMock)

	input := dto.CreateScheduleInput{
		RoomID:     "roomId",
		StartTime:  "17:00",
		EndTime:    "09:00",
		DaysOfWeek: []int{1, 3, 5},
	}

	_, err := uc.Execute(nil, input)

	require.ErrorIs(t, err, common.ErrInvalidScheduleTime)
}

func TestCreateSchedule_Execute_ScheduleAlreadyExists(t *testing.T) {
	storageMock := mocks.NewMockScheduleStorage(t)
	uc := NewCreateSchedule(storageMock)

	storageMock.EXPECT().IsScheduleExistsByRoomID(mock.Anything, "roomId").Return(true, nil)

	input := dto.CreateScheduleInput{
		RoomID:     "roomId",
		StartTime:  "09:00",
		EndTime:    "17:00",
		DaysOfWeek: []int{1, 3, 5},
	}

	_, err := uc.Execute(nil, input)

	require.ErrorIs(t, err, common.ErrScheduleExists)
}

func TestCreateSchedule_Execute_StorageError(t *testing.T) {
	storageMock := mocks.NewMockScheduleStorage(t)
	uc := NewCreateSchedule(storageMock)

	storageMock.EXPECT().IsScheduleExistsByRoomID(mock.Anything, "roomId").Return(false, assert.AnError)

	input := dto.CreateScheduleInput{
		RoomID:     "roomId",
		StartTime:  "09:00",
		EndTime:    "17:00",
		DaysOfWeek: []int{1, 3, 5},
	}

	_, err := uc.Execute(nil, input)

	require.ErrorIs(t, err, assert.AnError)
}
