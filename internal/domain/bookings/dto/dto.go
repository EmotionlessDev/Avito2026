package dto

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
