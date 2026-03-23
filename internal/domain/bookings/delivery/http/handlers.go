package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/usecases"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/helpers"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/middleware"
)

type CreateBookingUsecase interface {
	Execute(ctx context.Context, input usecases.CreateBookingInput) (*bookings.Booking, error)
}

type BookingHandler struct {
	createUsecase CreateBookingUsecase
}

func NewBookingHandler(uc CreateBookingUsecase) *BookingHandler {
	return &BookingHandler{
		createUsecase: uc,
	}
}

type createBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

type createBookingResponse struct {
	Booking *bookings.Booking `json:"booking"`
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.MethodNotAllowedResponse(w)
		return
	}

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "invalid_json",
		}, nil)
		return
	}

	if req.SlotID == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "slotId_required",
		}, nil)
		return
	}

	user, err := middleware.UserFromContext(r.Context())
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{
			"error": "unauthorized",
		}, nil)
		return
	}

	if user.Role != "user" {
		helpers.WriteJSON(w, http.StatusForbidden, helpers.Envelope{
			"error": "forbidden",
		}, nil)
		return
	}

	input := usecases.CreateBookingInput{
		SlotID:               req.SlotID,
		UserID:               user.UserID,
		CreateConferenceLink: req.CreateConferenceLink,
	}

	booking, err := h.createUsecase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidRequest):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
				"error": "invalid_request",
			}, nil)
		case errors.Is(err, common.ErrSlotNotFound):
			helpers.WriteJSON(w, http.StatusNotFound, helpers.Envelope{
				"error": "slot_not_found",
			}, nil)
		case errors.Is(err, common.ErrSlotAlreadyBooked):
			helpers.WriteJSON(w, http.StatusConflict, helpers.Envelope{
				"error": map[string]string{
					"code":    "SLOT_ALREADY_BOOKED",
					"message": "slot is already booked",
				},
			}, nil)
		default:
			helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
				"error": "internal_error",
			}, nil)
		}
		return
	}

	resp := createBookingResponse{
		Booking: booking,
	}

	helpers.WriteJSONObj(w, http.StatusCreated, resp, nil)
}
