package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
)

func TestCreateHandler_CreateRoom_Success(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	reqBody := CreateRoomRequest{
		Name:        "Test Room",
		Description: "Test description",
		Capacity:    10,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, "Test Room", "Test description", 10).
		Return("room-id-123", nil)

	handler.CreateRoom(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response CreateRoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "room-id-123", response.ID)
}

func TestCreateHandler_CreateRoom_WrongMethod(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/rooms/create", nil)
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCreateHandler_CreateRoom_InvalidJSON(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid_request", response["error"])
}

func TestCreateHandler_CreateRoom_NameRequired(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	reqBody := CreateRoomRequest{
		Name:        "",
		Description: "Test description",
		Capacity:    10,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "name_required", response["error"])
}

func TestCreateHandler_CreateRoom_InvalidCapacity(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	reqBody := CreateRoomRequest{
		Name:        "Test Room",
		Description: "Test description",
		Capacity:    -5,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid_capacity", response["error"])
}

func TestCreateHandler_CreateRoom_DuplicateRoom(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	reqBody := CreateRoomRequest{
		Name:        "Test Room",
		Description: "Test description",
		Capacity:    10,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, "Test Room", "Test description", 10).
		Return("", common.ErrDuplicateRoom)

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "room_exists", response["error"])
}

func TestCreateHandler_CreateRoom_InternalError(t *testing.T) {
	mockUsecase := mocks.NewMockCreateRoomUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	reqBody := CreateRoomRequest{
		Name:        "Test Room",
		Description: "Test description",
		Capacity:    10,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, "Test Room", "Test description", 10).
		Return("", errors.New("database error"))

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal_error", response["error"])
}
