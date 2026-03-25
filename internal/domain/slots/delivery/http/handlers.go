package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/helpers"
)

type GetSlotsUsecase interface {
	Execute(ctx context.Context, input dto.GetSlotsInput) ([]*slots.Slot, error)
}

type SlotHandler struct {
	getUsecase GetSlotsUsecase
}

func NewSlotHandler(usecase GetSlotsUsecase) *SlotHandler {
	return &SlotHandler{
		getUsecase: usecase,
	}
}

type GetSlotsResponse struct {
	Slots []*slots.Slot `json:"slots"`
}

func (h *SlotHandler) GetSlots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.MethodNotAllowedResponse(w)
		return
	}

	// path param
	roomID := r.PathValue("roomId")
	if roomID == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "room_id_required",
		}, nil)
		return
	}

	// query param
	date := r.URL.Query().Get("date")
	if date == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "date_required",
		}, nil)
		return
	}

	input := dto.GetSlotsInput{
		RoomID: roomID,
		Date:   date,
	}

	slotsList, err := h.getUsecase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidDate):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
				"error": "invalid_date",
			}, nil)
		case errors.Is(err, common.ErrRoomNotFound):
			helpers.WriteJSON(w, http.StatusNotFound, helpers.Envelope{
				"error": "room_not_found",
			}, nil)
		case errors.Is(err, common.ErrInvalidUUID):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
				"error": "invalid_room_id",
			}, nil)
		default:
			helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
				"error": "internal_error",
			}, nil)
		}
		return
	}

	resp := &GetSlotsResponse{
		Slots: slotsList,
	}

	helpers.WriteJSONObj(w, http.StatusOK, resp, nil)
}
