package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
)

type GetAllBookings struct {
	bookingStorage bookings.BookingStorage

	db *sql.DB
}

func NewGetAllBookings(bookingStorage bookings.BookingStorage, db *sql.DB) *GetAllBookings {
	return &GetAllBookings{
		bookingStorage: bookingStorage,
		db:             db,
	}
}

type GetAllBookingsInput struct {
	Page     int
	PageSize int
}

type Pagination struct {
	Page     int
	PageSize int
	Total    int
}

type GetAllBookingsOutput struct {
	Bookings   []*bookings.Booking
	Pagination Pagination
}

func (uc *GetAllBookings) Execute(ctx context.Context, input GetAllBookingsInput) (*GetAllBookingsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	offset := (input.Page - 1) * input.PageSize
	limit := input.PageSize

	// Start a transaction
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

	items, total, err := uc.bookingStorage.GetBookingsPaginated(ctx, tx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get bookings: %w", err)
	}
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, common.ErrCommitTx
	}
	committed = true

	return &GetAllBookingsOutput{
		Bookings: items,
		Pagination: Pagination{
			Page:     input.Page,
			PageSize: input.PageSize,
			Total:    total,
		}}, nil
}
