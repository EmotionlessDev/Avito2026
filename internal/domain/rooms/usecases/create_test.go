package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateRoom_Execute_Success(t *testing.T) {
	mockStorage := mocks.NewMockRoomStorage(t)
	uc := NewCreateRoom(mockStorage)

	mockStorage.EXPECT().CreateRoom(
		mock.Anything,
		"Test Room",
		"A room for testing",
		10,
	).Return("room-id-123", nil)

	id, err := uc.Execute(context.Background(), "Test Room", "A room for testing", 10)

	require.NoError(t, err)
	assert.Equal(t, "room-id-123", id)
}

func TestCreateRoom_Execute_Duplicate(t *testing.T) {
	mockStorage := mocks.NewMockRoomStorage(t)
	uc := NewCreateRoom(mockStorage)

	mockStorage.EXPECT().CreateRoom(
		mock.Anything,
		"Test Room",
		"A room for testing",
		10,
	).Return("", common.ErrDuplicateRoom)

	id, err := uc.Execute(context.Background(), "Test Room", "A room for testing", 10)

	require.ErrorIs(t, err, common.ErrDuplicateRoom)
	assert.Empty(t, id)
}

func TestCreateRoom_Execute_StorageError(t *testing.T) {
	mockStorage := mocks.NewMockRoomStorage(t)
	uc := NewCreateRoom(mockStorage)

	mockStorage.EXPECT().CreateRoom(
		mock.Anything,
		"Test Room",
		"A room for testing",
		10,
	).Return("", errors.New("database error"))

	id, err := uc.Execute(context.Background(), "Test Room", "A room for testing", 10)

	require.Error(t, err)
	assert.Empty(t, id)
	assert.Contains(t, err.Error(), "failed to create room")
}
