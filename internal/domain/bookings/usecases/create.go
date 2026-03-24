package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
)

type CreateBooking struct {
	bookingStorage bookings.BookingStorage
	slotStorage    slots.SlotStorage
	db             *sql.DB
}

func NewCreateBooking(bookingStorage bookings.BookingStorage, slotStorage slots.SlotStorage, db *sql.DB) *CreateBooking {
	return &CreateBooking{
		bookingStorage: bookingStorage,
		slotStorage:    slotStorage,
		db:             db,
	}
}

type CreateBookingInput struct {
	SlotID               string
	UserID               string
	CreateConferenceLink bool
}

func (uc *CreateBooking) Execute(ctx context.Context, input CreateBookingInput) (*bookings.Booking, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}

	tx, err := uc.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, common.ErrBeginTx
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// get slot
	slot, err := uc.slotStorage.GetSlotByID(ctx, tx, input.SlotID)
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
	booking, err := uc.bookingStorage.CreateBooking(ctx, tx, input.SlotID, input.UserID, link)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, common.ErrCommitTx
	}

	committed = true
	return booking, nil
}

func generateConferenceLink() string {
	// template
	return fmt.Sprintf("https://conference.example.com/%d", time.Now().UnixNano())
}
