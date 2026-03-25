package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/helpers"
)

type CreateScheduleUsecase interface {
	Execute(ctx context.Context, input dto.CreateScheduleInput) (*schedules.Schedule, error)
}

type ScheduleHandler struct {
	createUsecase CreateScheduleUsecase
}

func NewScheduleHandler(usecase CreateScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{
		createUsecase: usecase,
	}
}

type CreateScheduleRequest struct {
	StartTime  string `json:"start_time"` // "HH:MM"
	EndTime    string `json:"end_time"`   // "HH:MM"
	DaysOfWeek []int  `json:"days_of_week"`
}

type CreateScheduleResponse struct {
	Schedule *schedules.Schedule `json:"schedule"`
}

func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.MethodNotAllowedResponse(w)
		return
	}

	roomID := r.PathValue("roomId")
	if roomID == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "room_id_required",
		}, nil)
		return
	}

	var req CreateScheduleRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error":   "invalid_request",
			"message": err.Error(),
		}, nil)
		return
	}

	input := dto.CreateScheduleInput{
		RoomID:     roomID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		DaysOfWeek: req.DaysOfWeek,
	}

	schedule, err := h.createUsecase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidScheduleDay):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"error": "invalid_days_of_week"}, nil)
		case errors.Is(err, common.ErrDuplicateScheduleDay):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"error": "duplicate_days_of_week"}, nil)
		case errors.Is(err, common.ErrInvalidScheduleTime):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"error": "invalid_time"}, nil)
		case errors.Is(err, common.ErrScheduleExists):
			helpers.WriteJSON(w, http.StatusConflict, helpers.Envelope{
				"error": map[string]string{
					"code":    "SCHEDULE_EXISTS",
					"message": "schedule for this room already exists and cannot be changed",
				},
			}, nil)
		default:
			helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"error": "internal_error"}, nil)
		}
		return
	}

	rsp := &CreateScheduleResponse{Schedule: schedule}
	helpers.WriteJSONObj(w, http.StatusCreated, rsp, nil)
}
