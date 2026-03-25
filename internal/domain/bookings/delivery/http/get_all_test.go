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

func TestGetAllHandler_GetAllBookings_Success(t *testing.T) {
	mockUsecase := mocks.NewMockGetAllBookingsUsecase(t)
	handler := NewGetAllHandler(mockUsecase)

	userID := uuid.New().String()

	expectedBookings := []*bookings.Booking{
		{
			ID:             uuid.New().String(),
			SlotID:         uuid.New().String(),
			UserID:         userID,
			Status:         "confirmed",
			ConferenceLink: "link1",
			CreatedAt:      "2026-03-24T10:00:00Z",
		},
	}

	expectedPagination := dto.Pagination{
		Page:     1,
		PageSize: 20,
		Total:    1,
	}

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=1&pageSize=20", nil)

	claims := &middleware.Claims{
		UserID: userID,
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetAllBookingsInput{
			Page:     1,
			PageSize: 20,
		}).
		Return(&dto.GetAllBookingsOutput{
			Bookings:   expectedBookings,
			Pagination: expectedPagination,
		}, nil)

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp getAllBookingsResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp.Bookings, 1)
	assert.Equal(t, expectedBookings[0].ID, resp.Bookings[0].ID)
	assert.Equal(t, expectedPagination.Page, resp.Pagination.Page)
	assert.Equal(t, expectedPagination.Total, resp.Pagination.Total)
}

func TestGetAllHandler_GetAllBookings_DefaultPagination(t *testing.T) {
	mockUsecase := mocks.NewMockGetAllBookingsUsecase(t)
	handler := NewGetAllHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list", nil)

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.GetAllBookingsInput{
			Page:     1,
			PageSize: 20,
		}).
		Return(&dto.GetAllBookingsOutput{}, nil)

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllHandler_GetAllBookings_WrongMethod(t *testing.T) {
	handler := NewGetAllHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/bookings/list", nil)
	w := httptest.NewRecorder()

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestGetAllHandler_GetAllBookings_Unauthorized(t *testing.T) {
	handler := NewGetAllHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list", nil)
	w := httptest.NewRecorder()

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetAllHandler_GetAllBookings_Forbidden(t *testing.T) {
	handler := NewGetAllHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list", nil)

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "user",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetAllHandler_GetAllBookings_InvalidPage(t *testing.T) {
	handler := NewGetAllHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=abc", nil)

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllHandler_GetAllBookings_InvalidPageSize(t *testing.T) {
	handler := NewGetAllHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?pageSize=999", nil)

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllHandler_GetAllBookings_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockGetAllBookingsUsecase(t)
	handler := NewGetAllHandler(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list", nil)

	claims := &middleware.Claims{
		UserID: uuid.New().String(),
		Role:   "admin",
	}
	req = req.WithContext(middleware.WithUser(req.Context(), claims))

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, errors.New("some error"))

	handler.GetAllBookings(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
