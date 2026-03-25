package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/auth/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/helpers"
)

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type DummyLoginUsecase interface {
	Execute(ctx context.Context, role string) (dto.TokenResponse, error)
}

type Handler struct {
	dummyLoginUsecase DummyLoginUsecase
}

func NewHandler(dummyLogin DummyLoginUsecase) *Handler {
	return &Handler{
		dummyLoginUsecase: dummyLogin,
	}
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.WriteJSON(w, http.StatusMethodNotAllowed, helpers.Envelope{
			"error": "method_not_allowed",
		}, nil)
		return
	}

	var req DummyLoginRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error":   "invalid_request",
			"message": err.Error(),
		}, nil)
		return
	}

	ctx := r.Context()
	resp, err := h.dummyLoginUsecase.Execute(ctx, req.Role)
	if err != nil {
		if errors.Is(err, common.ErrInvalidRole) {
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
				"error":   "invalid_role",
				"message": "role must be 'admin' or 'user'",
			}, nil)
			return
		}

		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
			"error":   "internal_error",
			"message": err.Error(),
		}, nil)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{
		"token":      resp.Token,
		"user_id":    resp.UserID,
		"role":       resp.Role,
		"created_at": resp.CreatedAt,
	}, nil)

}
