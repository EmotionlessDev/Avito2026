package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/helpers"
)

type CreateRoomUsecase interface {
	Execute(ctx context.Context, name, description string, capacity int) (string, error)
}

type GetRoomsUsecase interface {
	Execute(ctx context.Context) ([]*rooms.Room, error)
}

type CreateHandler struct {
	createUsecase CreateRoomUsecase
}

type GetHandler struct {
	getUsecase GetRoomsUsecase
}

func NewCreateHandler(createUsecase CreateRoomUsecase) *CreateHandler {
	return &CreateHandler{
		createUsecase: createUsecase,
	}
}

func NewGetHandler(getUsecase GetRoomsUsecase) *GetHandler {
	return &GetHandler{
		getUsecase: getUsecase,
	}
}

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
}

type CreateRoomResponse struct {
	ID string `json:"id"`
}

func (h *CreateHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.MethodNotAllowedResponse(w)
		return
	}

	var req CreateRoomRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error":   "invalid_request",
			"message": err.Error(),
		}, nil)
		return
	}

	if req.Name == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "name_required",
		}, nil)
		return
	}

	if req.Capacity < 0 {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "invalid_capacity",
		}, nil)
		return

	}

	id, err := h.createUsecase.Execute(r.Context(), req.Name, req.Description, req.Capacity)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrDuplicateRoom):
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
				"error": "room_exists",
			}, nil)
		default:
			helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
				"error": "internal_error",
			}, nil)
		}
		return
	}

	rsp := &CreateRoomResponse{
		ID: id,
	}

	helpers.WriteJSONObj(w, http.StatusCreated, rsp, nil)

}

type GetRoomsResponse struct {
	Rooms []*rooms.Room `json:"rooms"`
}

func (h *GetHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.MethodNotAllowedResponse(w)
		return
	}

	roomList, err := h.getUsecase.Execute(r.Context())
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
			"error":   "internal_error",
			"message": err.Error(),
		}, nil)
		return
	}

	rsp := &GetRoomsResponse{
		Rooms: roomList,
	}

	helpers.WriteJSONObj(w, http.StatusOK, rsp, nil)
}
