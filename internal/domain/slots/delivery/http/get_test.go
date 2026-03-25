package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
)

func TestSlotHandler_GetSlots_Success(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	RoomID1 := uuid.New().String()
	RoomID2 := uuid.New().String()

	expectedSlots := []*slots.Slot{
		{
			ID:        "slot-1",
			RoomID:    RoomID1,
			StartTime: "2026-04-01T09:00:00Z",
			EndTime:   "2026-04-01T09:30:00Z",
		},
		{
			ID:        "slot-2",
			RoomID:    RoomID2,
			StartTime: "2026-04-01T09:30:00Z",
			EndTime:   "2026-04-01T10:00:00Z",
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+RoomID1+"/slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetSlotsInput{
			RoomID: RoomID1,
			Date:   "2026-04-01",
		}).
		Return(expectedSlots, nil)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response GetSlotsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Len(t, response.Slots, 2)
	require.Equal(t, "slot-1", response.Slots[0].ID)
	require.Equal(t, "slot-2", response.Slots[1].ID)
}

func TestSlotHandler_GetSlots_WrongMethod(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestSlotHandler_GetSlots_RoomIDRequired(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/rooms//slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.GetSlots(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "room_id_required", response["error"])
}

func TestSlotHandler_GetSlots_DateRequired(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/slots", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSlotHandler_GetSlots_InvalidDate(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/slots?date=invalid-date", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetSlotsInput{
			RoomID: roomID,
			Date:   "invalid-date",
		}).
		Return(nil, common.ErrInvalidDate)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "invalid_date", response["error"])
}

func TestSlotHandler_GetSlots_RoomNotFound(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetSlotsInput{
			RoomID: roomID,
			Date:   "2026-04-01",
		}).
		Return(nil, common.ErrRoomNotFound)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "room_not_found", response["error"])
}

func TestSlotHandler_GetSlots_InvalidUUID(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/invalid-uuid/slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetSlotsInput{
			RoomID: "invalid-uuid",
			Date:   "2026-04-01",
		}).
		Return(nil, common.ErrInvalidUUID)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "invalid_room_id", response["error"])
}

func TestSlotHandler_GetSlots_InternalError(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetSlotsInput{
			RoomID: roomID,
			Date:   "2026-04-01",
		}).
		Return(nil, common.ErrInvalidUUID)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "invalid_room_id", response["error"])
}

func TestSlotHandler_GetSlots_EmptyList(t *testing.T) {
	mockUsecase := mocks.NewMockGetSlotsUsecase(t)
	handler := NewSlotHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{roomId}/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/slots?date=2026-04-01", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetSlotsInput{
			RoomID: roomID,
			Date:   "2026-04-01",
		}).
		Return([]*slots.Slot{}, nil)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response GetSlotsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Empty(t, response.Slots)
}
