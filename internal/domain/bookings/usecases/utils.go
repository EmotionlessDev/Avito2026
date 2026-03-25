package usecases

import (
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots"
)

const (
	testBookingID = "550e8400-e29b-41d4-a716-446655440000"
	testUserID    = "660e8400-e29b-41d4-a716-446655440001"
	testSlotID    = "770e8400-e29b-41d4-a716-446655440002"
)

func createTestSlot(startTime time.Time) *slots.Slot {
	return &slots.Slot{
		ID:        testSlotID,
		RoomID:    "room-id",
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   startTime.Add(30 * time.Minute).Format(time.RFC3339),
	}
}

func createTestBooking(status string) *bookings.Booking {
	return &bookings.Booking{
		ID:             testBookingID,
		SlotID:         testSlotID,
		UserID:         testUserID,
		Status:         status,
		ConferenceLink: "https://meet.example.com/room123",
		CreatedAt:      time.Now().Format(time.RFC3339),
	}
}
