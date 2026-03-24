package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
)

type CancelBooking struct {
	bookingStorage bookings.BookingStorage

	db *sql.DB
}

func NewCancelBooking(bookingStorage bookings.BookingStorage, db *sql.DB) *CancelBooking {
	return &CancelBooking{
		bookingStorage: bookingStorage,
		db:             db,
	}
}

type CancelBookingInput struct {
	BookingID string
	UserID    string
}

type CancelBookingOutput struct {
	Booking *bookings.Booking
}

func (uc *CancelBooking) Execute(ctx context.Context, input CancelBookingInput) (*CancelBookingOutput, error) {
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

	// get booking
	booking, err := uc.bookingStorage.GetBookingByID(ctx, tx, input.BookingID)
	if err != nil {
		if errors.Is(err, common.ErrBookingNotFound) {
			return nil, common.ErrBookingNotFound
		}
		return nil, fmt.Errorf("get booking: %w", err)
	}

	// check if booking belongs to user
	if booking.UserID != input.UserID {
		return nil, common.ErrForbidden
	}

	// check if booking is already cancelled
	if booking.Status == "cancelled" {
		return &CancelBookingOutput{Booking: booking}, nil
	}

	// update booking status to cancelled
	err = uc.bookingStorage.UpdateBookingStatus(ctx, tx, input.BookingID, "cancelled")
	if err != nil {
		return nil, fmt.Errorf("update booking status: %w", err)
	}
	booking.Status = "cancelled" // update status in the returned booking object

	if err := tx.Commit(); err != nil {
		return nil, common.ErrCommitTx
	}
	committed = true

	return &CancelBookingOutput{Booking: booking}, nil
}
