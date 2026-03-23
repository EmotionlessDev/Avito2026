package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
)

type GetMyBookings struct {
	bookingStorage bookings.BookingStorage

	db *sql.DB
}

func NewGetMyBookings(bookingStorage bookings.BookingStorage, db *sql.DB) *GetMyBookings {
	return &GetMyBookings{
		bookingStorage: bookingStorage,
		db:             db,
	}
}

type GetMyBookingsInput struct {
	UserID string
}

type GetMyBookingsOutput struct {
	Bookings []*bookings.Booking
}

func (uc *GetMyBookings) Execute(ctx context.Context, input GetMyBookingsInput) (*GetMyBookingsOutput, error) {
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

	bookings, err := uc.bookingStorage.GetBookingsByUserID(ctx, tx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by user ID: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, common.ErrCommitTx
	}
	committed = true

	return &GetMyBookingsOutput{
		Bookings: bookings,
	}, nil
}
