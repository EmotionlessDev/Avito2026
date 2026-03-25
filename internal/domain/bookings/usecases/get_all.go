package usecases

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
)

type GetAllBookings struct {
	bookingStorage bookings.BookingStorage
}

func NewGetAllBookings(bookingStorage bookings.BookingStorage) *GetAllBookings {
	return &GetAllBookings{
		bookingStorage: bookingStorage,
	}
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

func (uc *GetAllBookings) Execute(ctx context.Context, input dto.GetAllBookingsInput) (*GetAllBookingsOutput, error) {
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

	items, total, err := uc.bookingStorage.GetBookingsPaginated(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get bookings: %w", err)
	}

	return &GetAllBookingsOutput{
		Bookings: items,
		Pagination: Pagination{
			Page:     input.Page,
			PageSize: input.PageSize,
			Total:    total,
		}}, nil
}
