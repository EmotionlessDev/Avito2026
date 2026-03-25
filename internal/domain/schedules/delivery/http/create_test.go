package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/middleware"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateScheduleHandler_Success(t *testing.T) {
	mockUsecase := mocks.NewMockCreateScheduleUsecase(t)
	handler := NewScheduleHandler(mockUsecase)

	roomID := uuid.New().String()

	expectedSchedule := &schedules.Schedule{
		ID:         uuid.New().String(),
		RoomID:     roomID,
		StartTime:  "09:00",
		EndTime:    "18:00",
		DaysOfWeek: []int{1, 2, 3},
	}

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
			middleware.RoleBased("admin"),
		),
	)

	body := `{
		"start_time": "09:00",
		"end_time": "18:00",
		"days_of_week": [1,2,3]
	}`

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/schedule/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	token := generateTestJWT(t, "test-secret", uuid.New().String(), "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, dto.CreateScheduleInput{
			RoomID:     roomID,
			StartTime:  "09:00",
			EndTime:    "18:00",
			DaysOfWeek: []int{1, 2, 3},
		}).
		Return(expectedSchedule, nil)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp CreateScheduleResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Schedule)
	assert.Equal(t, expectedSchedule.ID, resp.Schedule.ID)
}

func TestCreateScheduleHandler_InvalidDays(t *testing.T) {
	mockUsecase := mocks.NewMockCreateScheduleUsecase(t)
	handler := NewScheduleHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
			middleware.RoleBased("admin"),
		),
	)

	body := `{
		"start_time": "09:00",
		"end_time": "18:00",
		"days_of_week": [9]
	}`

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/schedule/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	token := generateTestJWT(t, "test-secret", uuid.New().String(), "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, common.ErrInvalidScheduleDay)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateScheduleHandler_DuplicateDays(t *testing.T) {
	mockUsecase := mocks.NewMockCreateScheduleUsecase(t)
	handler := NewScheduleHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
			middleware.RoleBased("admin"),
		),
	)

	body := `{
		"start_time": "09:00",
		"end_time": "18:00",
		"days_of_week": [1,1]
	}`

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/schedule/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	token := generateTestJWT(t, "test-secret", uuid.New().String(), "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, common.ErrDuplicateScheduleDay)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateScheduleHandler_InvalidTime(t *testing.T) {
	mockUsecase := mocks.NewMockCreateScheduleUsecase(t)
	handler := NewScheduleHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
			middleware.RoleBased("admin"),
		),
	)

	body := `{
		"start_time": "18:00",
		"end_time": "09:00",
		"days_of_week": [1]
	}`

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/schedule/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	token := generateTestJWT(t, "test-secret", uuid.New().String(), "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, common.ErrInvalidScheduleTime)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateScheduleHandler_ScheduleExists(t *testing.T) {
	mockUsecase := mocks.NewMockCreateScheduleUsecase(t)
	handler := NewScheduleHandler(mockUsecase)

	roomID := uuid.New().String()

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
			middleware.RoleBased("admin"),
		),
	)

	body := `{
		"start_time": "09:00",
		"end_time": "18:00",
		"days_of_week": [1]
	}`

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/schedule/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	token := generateTestJWT(t, "test-secret", uuid.New().String(), "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, mock.Anything).
		Return(nil, common.ErrScheduleExists)

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateScheduleHandler_Unauthorized(t *testing.T) {
	handler := NewScheduleHandler(nil)

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
		),
	)

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+uuid.New().String()+"/schedule/create", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateScheduleHandler_Forbidden(t *testing.T) {
	handler := NewScheduleHandler(nil)

	mux := http.NewServeMux()
	mux.Handle("/rooms/{roomId}/schedule/create",
		middleware.Chain(
			http.HandlerFunc(handler.CreateSchedule),
			middleware.JWTMiddleware("test-secret"),
			middleware.RoleBased("admin"),
		),
	)

	token := generateTestJWT(t, "test-secret", uuid.New().String(), "user")

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+uuid.New().String()+"/schedule/create", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func generateTestJWT(t *testing.T, secret, userID, role string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	return str
}
