package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/auth/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_DummyLogin_Success_User(t *testing.T) {
	mockUsecase := mocks.NewMockDummyLoginUsecase(t)
	handler := &Handler{
		dummyLoginUsecase: mockUsecase,
	}

	body := `{"role":"user"}`
	req := httptest.NewRequest(http.MethodPost, "/dummylogin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	expectedResp := dto.TokenResponse{
		Token:     "test-token",
		UserID:    uuid.New().String(),
		Role:      "user",
		CreatedAt: time.Now().UTC(),
	}

	mockUsecase.EXPECT().
		Execute(mock.Anything, "user").
		Return(expectedResp, nil)

	handler.DummyLogin(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, expectedResp.Token, resp["token"])
	assert.Equal(t, expectedResp.UserID, resp["user_id"])
	assert.Equal(t, expectedResp.Role, resp["role"])
}

func TestHandler_DummyLogin_Success_Admin(t *testing.T) {
	mockUsecase := mocks.NewMockDummyLoginUsecase(t)
	handler := &Handler{
		dummyLoginUsecase: mockUsecase,
	}

	body := `{"role":"admin"}`
	req := httptest.NewRequest(http.MethodPost, "/dummylogin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	expectedResp := dto.TokenResponse{
		Token:     "test-token",
		UserID:    uuid.New().String(),
		Role:      "admin",
		CreatedAt: time.Now().UTC(),
	}

	mockUsecase.EXPECT().
		Execute(mock.Anything, "admin").
		Return(expectedResp, nil)

	handler.DummyLogin(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, expectedResp.Token, resp["token"])
	assert.Equal(t, expectedResp.UserID, resp["user_id"])
	assert.Equal(t, expectedResp.Role, resp["role"])
}

func TestHandler_DummyLogin_InvalidRole(t *testing.T) {
	mockUsecase := mocks.NewMockDummyLoginUsecase(t)
	handler := &Handler{
		dummyLoginUsecase: mockUsecase,
	}

	body := `{"role":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/dummylogin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, "invalid").
		Return(dto.TokenResponse{}, common.ErrInvalidRole)

	handler.DummyLogin(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, "invalid_role", resp["error"])
}

func TestHandler_DummyLogin_InternalError(t *testing.T) {
	mockUsecase := mocks.NewMockDummyLoginUsecase(t)
	handler := &Handler{
		dummyLoginUsecase: mockUsecase,
	}

	body := `{"role":"user"}`
	req := httptest.NewRequest(http.MethodPost, "/dummylogin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mockUsecase.EXPECT().
		Execute(mock.Anything, "user").
		Return(dto.TokenResponse{}, errors.New("some error"))

	handler.DummyLogin(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, "internal_error", resp["error"])
}

func TestHandler_DummyLogin_MethodNotAllowed(t *testing.T) {
	mockUsecase := mocks.NewMockDummyLoginUsecase(t)
	handler := &Handler{
		dummyLoginUsecase: mockUsecase,
	}

	req := httptest.NewRequest(http.MethodGet, "/dummylogin", nil)

	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	require.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, "method_not_allowed", resp["error"])
}
