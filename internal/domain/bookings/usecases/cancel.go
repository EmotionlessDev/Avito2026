package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
)

type CancelBooking struct {
	bookingStorage bookings.BookingStorage
}

func NewCancelBooking(bookingStorage bookings.BookingStorage) *CancelBooking {
	return &CancelBooking{
		bookingStorage: bookingStorage,
	}
}

type CancelBookingOutput struct {
	Booking *bookings.Booking
}

func (uc *CancelBooking) Execute(ctx context.Context, input dto.CancelBookingInput) (*CancelBookingOutput, error) {
	// get booking
	booking, err := uc.bookingStorage.GetBookingByID(ctx, input.BookingID)
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
	err = uc.bookingStorage.UpdateBookingStatus(ctx, input.BookingID, "cancelled")
	if err != nil {
		return nil, fmt.Errorf("update booking status: %w", err)
	}
	booking.Status = "cancelled" // update status in the returned booking object

	return &CancelBookingOutput{Booking: booking}, nil
}
