package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMyBookings_Execute_Success(t *testing.T) {
	// Arrange
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetMyBookings(bookingStorageMock)

	testBookings := []*bookings.Booking{
		createTestBooking("active"),
	}

	input := dto.GetMyBookingsInput{
		UserID: testUserID,
	}

	bookingStorageMock.EXPECT().GetBookingsByUserID(mock.Anything, testUserID).Return(testBookings, nil)
	// Act

	output, err := uc.Execute(context.Background(), input)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, testBookings, output.Bookings)
}

func TestGetMyBookings_Execute_Success_EmptyList(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetMyBookings(bookingStorageMock)

	input := dto.GetMyBookingsInput{
		UserID: testUserID,
	}

	bookingStorageMock.EXPECT().
		GetBookingsByUserID(mock.Anything, testUserID).
		Return([]*bookings.Booking{}, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Empty(t, output.Bookings)
	assert.Len(t, output.Bookings, 0)
}

func TestGetMyBookings_Execute_StorageError(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetMyBookings(bookingStorageMock)

	input := dto.GetMyBookingsInput{
		UserID: testUserID,
	}

	bookingStorageMock.EXPECT().
		GetBookingsByUserID(mock.Anything, testUserID).
		Return(nil, errors.New("database connection failed"))

	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "get bookings by user ID")
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestGetMyBookings_Execute_MixedStatuses(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetMyBookings(bookingStorageMock)

	testBookings := []*bookings.Booking{
		{
			ID:             "booking-1",
			SlotID:         testSlotID,
			UserID:         testUserID,
			Status:         "active",
			ConferenceLink: "https://meet.example.com/room1",
			CreatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             "booking-2",
			SlotID:         testSlotID,
			UserID:         testUserID,
			Status:         "cancelled",
			ConferenceLink: "https://meet.example.com/room2",
			CreatedAt:      time.Now().Format(time.RFC3339),
		},
		{
			ID:             "booking-3",
			SlotID:         testSlotID,
			UserID:         testUserID,
			Status:         "completed",
			ConferenceLink: "",
			CreatedAt:      time.Now().Format(time.RFC3339),
		},
	}

	input := dto.GetMyBookingsInput{
		UserID: testUserID,
	}

	bookingStorageMock.EXPECT().
		GetBookingsByUserID(mock.Anything, testUserID).
		Return(testBookings, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Bookings, 3)
	assert.Equal(t, "active", output.Bookings[0].Status)
	assert.Equal(t, "cancelled", output.Bookings[1].Status)
	assert.Equal(t, "completed", output.Bookings[2].Status)
}

func TestGetMyBookings_Execute_CorrectUserID(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetMyBookings(bookingStorageMock)

	customUserID := "custom-user-id-123"

	input := dto.GetMyBookingsInput{
		UserID: customUserID,
	}

	bookingStorageMock.EXPECT().
		GetBookingsByUserID(mock.Anything, customUserID).
		Return([]*bookings.Booking{}, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
}
