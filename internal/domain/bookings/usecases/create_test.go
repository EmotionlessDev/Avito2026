package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
)

func TestCreateBooking_Execute_Success_NoConferenceLink(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	uc := NewCreateBooking(bookingStorageMock, slotStorageMock)

	futureTime := time.Now().UTC().Add(24 * time.Hour)
	testSlot := createTestSlot(futureTime)

	slotStorageMock.EXPECT().
		GetSlotByID(mock.Anything, testSlotID).
		Return(testSlot, nil)

	bookingStorageMock.EXPECT().
		CreateBooking(mock.Anything, testSlotID, testUserID, (*string)(nil)).
		Return(createTestBooking("active"), nil)

	input := dto.CreateBookingInput{
		SlotID:               testSlotID,
		UserID:               testUserID,
		CreateConferenceLink: false,
	}

	booking, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, booking)
	assert.Equal(t, testBookingID, booking.ID)
	assert.Equal(t, testSlotID, booking.SlotID)
	assert.Equal(t, testUserID, booking.UserID)
}

func TestCreateBooking_Execute_Success_WithConferenceLink(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	uc := NewCreateBooking(bookingStorageMock, slotStorageMock)

	futureTime := time.Now().UTC().Add(24 * time.Hour)
	testSlot := createTestSlot(futureTime)

	slotStorageMock.EXPECT().
		GetSlotByID(mock.Anything, testSlotID).
		Return(testSlot, nil)

	bookingStorageMock.EXPECT().
		CreateBooking(mock.Anything, testSlotID, testUserID, mock.MatchedBy(func(link *string) bool {
			return link != nil && len(*link) > 0
		})).
		Return(createTestBooking("active"), nil)

	input := dto.CreateBookingInput{
		SlotID:               testSlotID,
		UserID:               testUserID,
		CreateConferenceLink: true,
	}

	booking, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, booking)
	assert.Equal(t, testBookingID, booking.ID)
}

func TestCreateBooking_Execute_SlotNotFound(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	uc := NewCreateBooking(bookingStorageMock, slotStorageMock)

	slotStorageMock.EXPECT().
		GetSlotByID(mock.Anything, testSlotID).
		Return(nil, common.ErrSlotNotFound)

	input := dto.CreateBookingInput{
		SlotID:               testSlotID,
		UserID:               testUserID,
		CreateConferenceLink: false,
	}

	booking, err := uc.Execute(context.Background(), input)

	assert.ErrorIs(t, err, common.ErrSlotNotFound)
	assert.Nil(t, booking)
}

func TestCreateBooking_Execute_SlotInPast(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	uc := NewCreateBooking(bookingStorageMock, slotStorageMock)

	pastTime := time.Now().UTC().Add(-24 * time.Hour)
	testSlot := createTestSlot(pastTime)

	slotStorageMock.EXPECT().
		GetSlotByID(mock.Anything, testSlotID).
		Return(testSlot, nil)

	input := dto.CreateBookingInput{
		SlotID:               testSlotID,
		UserID:               testUserID,
		CreateConferenceLink: false,
	}

	booking, err := uc.Execute(context.Background(), input)

	assert.ErrorIs(t, err, common.ErrInvalidRequest)
	assert.Nil(t, booking)
}

func TestCreateBooking_Execute_CreateBookingError(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	uc := NewCreateBooking(bookingStorageMock, slotStorageMock)

	futureTime := time.Now().UTC().Add(24 * time.Hour)
	testSlot := createTestSlot(futureTime)

	slotStorageMock.EXPECT().
		GetSlotByID(mock.Anything, testSlotID).
		Return(testSlot, nil)

	bookingStorageMock.EXPECT().
		CreateBooking(mock.Anything, testSlotID, testUserID, (*string)(nil)).
		Return(nil, common.ErrSlotAlreadyBooked)

	input := dto.CreateBookingInput{
		SlotID:               testSlotID,
		UserID:               testUserID,
		CreateConferenceLink: false,
	}

	booking, err := uc.Execute(context.Background(), input)

	assert.ErrorIs(t, err, common.ErrSlotAlreadyBooked)
	assert.Nil(t, booking)
}

func TestCreateBooking_Execute_GetSlotError(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	slotStorageMock := mocks.NewMockSlotStorage(t)
	uc := NewCreateBooking(bookingStorageMock, slotStorageMock)

	slotStorageMock.EXPECT().
		GetSlotByID(mock.Anything, testSlotID).
		Return(nil, errors.New("database connection failed"))

	input := dto.CreateBookingInput{
		SlotID:               testSlotID,
		UserID:               testUserID,
		CreateConferenceLink: false,
	}

	booking, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, booking)
	assert.Contains(t, err.Error(), "database connection failed")
}
