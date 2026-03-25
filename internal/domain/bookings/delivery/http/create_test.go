package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/middleware"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateHandler_CreateBooking_Success(t *testing.T) {
	mockUsecase := mocks.NewMockCreateBookingUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	slotID := uuid.New().String()
	userID := uuid.New().String()

	expectedBooking := &bookings.Booking{
		ID:             uuid.New().String(),
		SlotID:         slotID,
		UserID:         userID,
		Status:         "confirmed",
		ConferenceLink: "https://meet.example.com/test",
		CreatedAt:      "2026-03-24T10:00:00Z",
	}

	body := `{
		"slotId": "` + slotID + `",
		"createConferenceLink": true
	}`

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.CreateBookingInput{
			SlotID:               slotID,
			UserID:               userID,
			CreateConferenceLink: true,
		}).
		Return(expectedBooking, nil)

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp createBookingResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, expectedBooking.ID, resp.Booking.ID)
	assert.Equal(t, expectedBooking.SlotID, resp.Booking.SlotID)
	assert.Equal(t, expectedBooking.UserID, resp.Booking.UserID)
}

func TestCreateHandler_CreateBooking_WrongMethod(t *testing.T) {
	handler := NewCreateHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/create", nil)
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCreateHandler_CreateBooking_InvalidJSON(t *testing.T) {
	handler := NewCreateHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader("{invalid"))
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateHandler_CreateBooking_MissingSlotID(t *testing.T) {
	handler := NewCreateHandler(nil)

	body := `{}`

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateHandler_CreateBooking_Unauthorized(t *testing.T) {
	handler := NewCreateHandler(nil)

	body := `{"slotId":"` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateHandler_CreateBooking_Forbidden(t *testing.T) {
	handler := NewCreateHandler(nil)

	body := `{"slotId":"` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(body))

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateHandler_CreateBooking_SlotNotFound(t *testing.T) {
	mockUsecase := mocks.NewMockCreateBookingUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	slotID := uuid.New().String()
	userID := uuid.New().String()

	body := `{"slotId":"` + slotID + `"}`

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(body))

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, common.ErrSlotNotFound)

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateHandler_CreateBooking_AlreadyBooked(t *testing.T) {
	mockUsecase := mocks.NewMockCreateBookingUsecase(t)
	handler := NewCreateHandler(mockUsecase)

	slotID := uuid.New().String()
	userID := uuid.New().String()

	body := `{"slotId":"` + slotID + `"}`

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(body))

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, common.ErrSlotAlreadyBooked)

	handler.CreateBooking(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}
