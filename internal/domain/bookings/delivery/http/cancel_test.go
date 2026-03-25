package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/middleware"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCancelHandler_CancelBooking_Success(t *testing.T) {
	mockUsecase := mocks.NewMockCancelBookingUsecase(t)
	handler := NewCancelHandler(mockUsecase)

	bookingID := uuid.New().String()
	userID := uuid.New().String()

	expectedBooking := &bookings.Booking{
		ID:             bookingID,
		SlotID:         uuid.New().String(),
		UserID:         userID,
		Status:         "cancelled",
		ConferenceLink: "https://meet.example.com/room123",
		CreatedAt:      "2026-03-24T10:00:00Z",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/bookings/{bookingId}/cancel", handler.CancelBooking)

	req := httptest.NewRequest(http.MethodPost, "/bookings/"+bookingID+"/cancel", nil)
	req.Header.Set("Content-Type", "application/json")

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	ctx := middleware.WithUser(req.Context(), claims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.CancelBookingInput{
			BookingID: bookingID,
			UserID:    userID,
		}).
		Return(&dto.CancelBookingOutput{Booking: expectedBooking}, nil)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response CancelBookingResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	require.NotNil(t, response.Booking)
	assert.Equal(t, expectedBooking.ID, response.Booking.ID)
	assert.Equal(t, expectedBooking.SlotID, response.Booking.SlotID)
	assert.Equal(t, expectedBooking.UserID, response.Booking.UserID)
	assert.Equal(t, expectedBooking.Status, response.Booking.Status)
	assert.Equal(t, expectedBooking.ConferenceLink, response.Booking.ConferenceLink)
	assert.Equal(t, expectedBooking.CreatedAt, response.Booking.CreatedAt)
}

func TestCancelHandler_CancelBooking_WrongMethod(t *testing.T) {
	mockUsecase := mocks.NewMockCancelBookingUsecase(t)
	handler := NewCancelHandler(mockUsecase)

	bookingID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/bookings/{bookingId}/cancel", handler.CancelBooking)

	req := httptest.NewRequest(http.MethodGet, "/bookings/"+bookingID+"/cancel", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCancelHandler_CancelBooking_InvalidBookingID(t *testing.T) {
	mockUsecase := mocks.NewMockCancelBookingUsecase(t)
	handler := NewCancelHandler(mockUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/bookings/{bookingId}/cancel", handler.CancelBooking)

	req := httptest.NewRequest(http.MethodPost, "/bookings/invalid-booking-id/cancel", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestCancelHandler_CancelBooking_Unauthorized(t *testing.T) {
	mockUsecase := mocks.NewMockCancelBookingUsecase(t)
	handler := NewCancelHandler(mockUsecase)

	bookingID := uuid.New().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/bookings/{bookingId}/cancel", handler.CancelBooking)

	req := httptest.NewRequest(http.MethodPost, "/bookings/"+bookingID+"/cancel", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
