package dto

import "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"

type CancelBookingInput struct {
	BookingID string
	UserID    string
}

type CreateBookingInput struct {
	SlotID               string
	UserID               string
	CreateConferenceLink bool
}

type GetAllBookingsInput struct {
	Page     int
	PageSize int
}

type GetMyBookingsInput struct {
	UserID string
}

type CancelBookingOutput struct {
	Booking *bookings.Booking
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

type GetMyBookingsOutput struct {
	Bookings []*bookings.Booking
}
