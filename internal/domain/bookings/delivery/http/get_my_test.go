package http

import (
	"encoding/json"
	"errors"
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

func TestGetMyHandler_GetMyBookings_Success(t *testing.T) {
	mockUsecase := mocks.NewMockGetMyBookingsUsecase(t)
	handler := NewGetMyHandler(mockUsecase)

	userID := uuid.New().String()

	expectedBookings := []*bookings.Booking{
		{
			ID:             uuid.New().String(),
			SlotID:         uuid.New().String(),
			UserID:         userID,
			Status:         "confirmed",
			ConferenceLink: "https://meet.example.com/room1",
			CreatedAt:      "2026-03-24T10:00:00Z",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetMyBookingsInput{
			UserID: userID,
		}).
		Return(&dto.GetMyBookingsOutput{
			Bookings: expectedBookings,
		}, nil)

	handler.GetMyBookings(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp getMyBookingsResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp.Bookings, 1)
	assert.Equal(t, expectedBookings[0].ID, resp.Bookings[0].ID)
	assert.Equal(t, expectedBookings[0].UserID, resp.Bookings[0].UserID)
}

func TestGetMyHandler_GetMyBookings_EmptyList(t *testing.T) {
	mockUsecase := mocks.NewMockGetMyBookingsUsecase(t)
	handler := NewGetMyHandler(mockUsecase)

	userID := uuid.New().String()

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetMyBookingsInput{
			UserID: userID,
		}).
		Return(&dto.GetMyBookingsOutput{
			Bookings: []*bookings.Booking{},
		}, nil)

	handler.GetMyBookings(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp getMyBookingsResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp.Bookings, 0)
}

func TestGetMyHandler_GetMyBookings_WrongMethod(t *testing.T) {
	handler := NewGetMyHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/bookings/my", nil)
	w := httptest.NewRecorder()

	handler.GetMyBookings(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestGetMyHandler_GetMyBookings_Unauthorized(t *testing.T) {
	handler := NewGetMyHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
	w := httptest.NewRecorder()

	handler.GetMyBookings(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMyHandler_GetMyBookings_Forbidden(t *testing.T) {
	handler := NewGetMyHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	handler.GetMyBookings(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetMyHandler_GetMyBookings_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockGetMyBookingsUsecase(t)
	handler := NewGetMyHandler(mockUsecase)

	userID := uuid.New().String()

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetMyBookingsInput{
			UserID: userID,
		}).
		Return(nil, errors.New("some error"))

	handler.GetMyBookings(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
