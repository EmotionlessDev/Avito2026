package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetHandler_GetRooms_Success(t *testing.T) {
	mockUsecase := mocks.NewMockGetRoomsUsecase(t)
	handler := NewGetHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().Execute(mock.Anything).Return(
		[]*rooms.Room{
			{
				ID:          "room-id-123",
				Name:        "Test Room",
				Description: "Test description",
				Capacity:    10,
			},
			{
				ID:          "room-id-456",
				Name:        "Another Room",
				Description: "Another description",
				Capacity:    20,
			},
		}, nil,
	)

	handler.GetRooms(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response GetRoomsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Len(t, response.Rooms, 2)
	require.Equal(t, "room-id-123", response.Rooms[0].ID)
	require.Equal(t, "Test Room", response.Rooms[0].Name)
	require.Equal(t, "Test description", response.Rooms[0].Description)
	require.Equal(t, 10, response.Rooms[0].Capacity)
	require.Equal(t, "room-id-456", response.Rooms[1].ID)
	require.Equal(t, "Another Room", response.Rooms[1].Name)
	require.Equal(t, "Another description", response.Rooms[1].Description)
	require.Equal(t, 20, response.Rooms[1].Capacity)
}

func TestGetHandler_GetRooms_WrongMethod(t *testing.T) {
	mockUsecase := mocks.NewMockGetRoomsUsecase(t)
	handler := NewGetHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodPost, "/rooms/list", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.GetRooms(w, req)
	require.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response GetRoomsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Empty(t, response.Rooms)
}

func TestGetHandler_GetRooms_InternalError(t *testing.T) {
	mockUsecase := mocks.NewMockGetRoomsUsecase(t)
	handler := NewGetHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().Execute(mock.Anything).Return(
		[]*rooms.Room{
			{
				ID:          "room-id-123",
				Name:        "Test Room",
				Description: "Test description",
				Capacity:    10,
			},
			{
				ID:          "room-id-456",
				Name:        "Another Room",
				Description: "Another description",
				Capacity:    20,
			},
		}, common.ErrInvalidUUID,
	)

	handler.GetRooms(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var response GetRoomsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Empty(t, response.Rooms)
}
