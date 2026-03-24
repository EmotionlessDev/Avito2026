package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// type GetRooms struct {
// 	roomStorage rooms.RoomStorage
// }
//
// func NewGetRooms(roomStorage rooms.RoomStorage) *GetRooms {
// 	return &GetRooms{
// 		roomStorage: roomStorage,
// 	}
// }
//
// func (uc *GetRooms) Execute(ctx context.Context) ([]*rooms.Room, error) {
// 	roomList, err := uc.roomStorage.GetRooms(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get rooms: %w", err)
// 	}
//
// 	return roomList, nil
// }

func TestGetRooms_Execute_Success(t *testing.T) {
	mockStorage := mocks.NewMockRoomStorage(t)
	uc := NewGetRooms(mockStorage)

	expectedRooms := []*rooms.Room{
		{ID: "1", Name: "Room 1", Description: "Description 1", Capacity: 10},
		{ID: "2", Name: "Room 2", Description: "Description 2", Capacity: 20},
	}

	mockStorage.EXPECT().GetRooms(mock.Anything).Return(expectedRooms, nil)

	actualRooms, err := uc.Execute(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedRooms, actualRooms)
}

func TestGetRooms_Execute_EmptyList(t *testing.T) {
	mockStorage := mocks.NewMockRoomStorage(t)
	uc := NewGetRooms(mockStorage)

	mockStorage.EXPECT().GetRooms(mock.Anything).Return([]*rooms.Room{}, nil)

	actualRooms, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Empty(t, actualRooms)
}

func TestGetRooms_Execute_Error(t *testing.T) {
	mockStorage := mocks.NewMockRoomStorage(t)
	uc := NewGetRooms(mockStorage)

	mockStorage.EXPECT().GetRooms(mock.Anything).Return(nil, errors.New("database error"))

	actualRooms, err := uc.Execute(context.Background())

	assert.Error(t, err)
	assert.Nil(t, actualRooms)
	assert.Contains(t, err.Error(), "failed to get rooms: database error")
}
