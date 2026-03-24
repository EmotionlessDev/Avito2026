package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
)

type CreateBooking struct {
	bookingStorage bookings.BookingStorage
	slotStorage    slots.SlotStorage
}

func NewCreateBooking(bookingStorage bookings.BookingStorage, slotStorage slots.SlotStorage) *CreateBooking {
	return &CreateBooking{
		bookingStorage: bookingStorage,
		slotStorage:    slotStorage,
	}
}

type CreateBookingInput struct {
	SlotID               string
	UserID               string
	CreateConferenceLink bool
}

func (uc *CreateBooking) Execute(ctx context.Context, input CreateBookingInput) (*bookings.Booking, error) {
	// get slot
	slot, err := uc.slotStorage.GetSlotByID(ctx, input.SlotID)
	if err != nil {
		return nil, err
	}

	// check if slot is in the past
	startTime, err := time.Parse(time.RFC3339, slot.StartTime)
	if err != nil {
		return nil, fmt.Errorf("parse slot time: %w", err)
	}

	if startTime.Before(time.Now().UTC()) {
		return nil, common.ErrInvalidRequest
	}

	// link generation
	var link *string
	if input.CreateConferenceLink {
		url := generateConferenceLink()
		link = &url
	}

	// create booking
	booking, err := uc.bookingStorage.CreateBooking(ctx, input.SlotID, input.UserID, link)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func generateConferenceLink() string {
	// template
	return fmt.Sprintf("https://conference.example.com/%d", time.Now().UnixNano())
}
