package usecases

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
)

type GetMyBookings struct {
	bookingStorage bookings.BookingStorage
}

func NewGetMyBookings(bookingStorage bookings.BookingStorage) *GetMyBookings {
	return &GetMyBookings{
		bookingStorage: bookingStorage,
	}
}

type GetMyBookingsInput struct {
	UserID string
}

type GetMyBookingsOutput struct {
	Bookings []*bookings.Booking
}

func (uc *GetMyBookings) Execute(ctx context.Context, input GetMyBookingsInput) (*GetMyBookingsOutput, error) {
	bookings, err := uc.bookingStorage.GetBookingsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by user ID: %w", err)
	}

	return &GetMyBookingsOutput{
		Bookings: bookings,
	}, nil
}
