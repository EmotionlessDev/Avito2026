package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCancelBooking_Execute_Success(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewCancelBooking(bookingStorageMock)

	existingBooking := createTestBooking("active")

	bookingStorageMock.EXPECT().
		GetBookingByID(mock.Anything, testBookingID).
		Return(existingBooking, nil)

	bookingStorageMock.EXPECT().
		UpdateBookingStatus(mock.Anything, testBookingID, "cancelled").
		Return(nil)

	input := dto.CancelBookingInput{
		BookingID: testBookingID,
		UserID:    testUserID,
	}

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	require.NotNil(t, output.Booking)
	assert.Equal(t, testBookingID, output.Booking.ID)
	assert.Equal(t, "cancelled", output.Booking.Status)
	assert.Equal(t, testUserID, output.Booking.UserID)
}

func TestCancelBooking_Execute_BookingNotFound(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewCancelBooking(bookingStorageMock)

	bookingStorageMock.EXPECT().
		GetBookingByID(mock.Anything, testBookingID).
		Return(nil, common.ErrBookingNotFound)

	input := dto.CancelBookingInput{
		BookingID: testBookingID,
		UserID:    testUserID,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.ErrorIs(t, err, common.ErrBookingNotFound)
	assert.Nil(t, output)
}

func TestCancelBooking_Execute_Forbidden(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewCancelBooking(bookingStorageMock)

	existingBooking := createTestBooking("active")
	existingBooking.UserID = "other-user-id"

	bookingStorageMock.EXPECT().
		GetBookingByID(mock.Anything, testBookingID).
		Return(existingBooking, nil)

	input := dto.CancelBookingInput{
		BookingID: testBookingID,
		UserID:    testUserID,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.ErrorIs(t, err, common.ErrForbidden)
	assert.Nil(t, output)
}

func TestCancelBooking_Execute_AlreadyCancelled(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewCancelBooking(bookingStorageMock)

	existingBooking := createTestBooking("cancelled")

	bookingStorageMock.EXPECT().
		GetBookingByID(mock.Anything, testBookingID).
		Return(existingBooking, nil)

	input := dto.CancelBookingInput{
		BookingID: testBookingID,
		UserID:    testUserID,
	}

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	require.NotNil(t, output.Booking)
	assert.Equal(t, "cancelled", output.Booking.Status)
}

func TestCancelBooking_Execute_UpdateStatusError(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewCancelBooking(bookingStorageMock)

	existingBooking := createTestBooking("active")

	bookingStorageMock.EXPECT().
		GetBookingByID(mock.Anything, testBookingID).
		Return(existingBooking, nil)

	bookingStorageMock.EXPECT().
		UpdateBookingStatus(mock.Anything, testBookingID, "cancelled").
		Return(errors.New("database error"))

	input := dto.CancelBookingInput{
		BookingID: testBookingID,
		UserID:    testUserID,
	}

	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "update booking status")
	assert.Contains(t, err.Error(), "database error")
}

func TestCancelBooking_Execute_GetBookingError(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewCancelBooking(bookingStorageMock)

	bookingStorageMock.EXPECT().
		GetBookingByID(mock.Anything, testBookingID).
		Return(nil, errors.New("database connection failed"))

	input := dto.CancelBookingInput{
		BookingID: testBookingID,
		UserID:    testUserID,
	}

	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "get booking")
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestCancelBooking_Execute(t *testing.T) {
	tests := []struct {
		name            string
		bookingStatus   string
		bookingUserID   string
		inputUserID     string
		getBookingErr   error
		updateStatusErr error
		wantErr         error
		wantOutput      bool
	}{
		{
			name:          "success",
			bookingStatus: "active",
			bookingUserID: testUserID,
			inputUserID:   testUserID,
			wantErr:       nil,
			wantOutput:    true,
		},
		{
			name:          "booking not found",
			bookingStatus: "",
			bookingUserID: "",
			inputUserID:   testUserID,
			getBookingErr: common.ErrBookingNotFound,
			wantErr:       common.ErrBookingNotFound,
			wantOutput:    false,
		},
		{
			name:          "forbidden",
			bookingStatus: "active",
			bookingUserID: "other-user-id",
			inputUserID:   testUserID,
			wantErr:       common.ErrForbidden,
			wantOutput:    false,
		},
		{
			name:          "already cancelled",
			bookingStatus: "cancelled",
			bookingUserID: testUserID,
			inputUserID:   testUserID,
			wantErr:       nil,
			wantOutput:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookingStorageMock := mocks.NewMockBookingStorage(t)
			uc := NewCancelBooking(bookingStorageMock)

			if tt.getBookingErr != nil {
				bookingStorageMock.EXPECT().
					GetBookingByID(mock.Anything, testBookingID).
					Return(nil, tt.getBookingErr)
			} else {
				booking := createTestBooking(tt.bookingStatus)
				booking.UserID = tt.bookingUserID
				bookingStorageMock.EXPECT().
					GetBookingByID(mock.Anything, testBookingID).
					Return(booking, nil)

				if tt.wantErr == nil && tt.bookingStatus != "cancelled" {
					bookingStorageMock.EXPECT().
						UpdateBookingStatus(mock.Anything, testBookingID, "cancelled").
						Return(tt.updateStatusErr)
				}
			}

			input := dto.CancelBookingInput{
				BookingID: testBookingID,
				UserID:    tt.inputUserID,
			}

			output, err := uc.Execute(context.Background(), input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			if tt.wantOutput {
				assert.NotNil(t, output)
			} else {
				assert.Nil(t, output)
			}
		})
	}
}
