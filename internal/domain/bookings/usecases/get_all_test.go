package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/mocks"
)

func TestGetAllBookings_Execute_Success_Page1(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	testBookings := []*bookings.Booking{
		createTestBooking("active"),
		createTestBooking("active"),
	}
	totalCount := 50

	input := GetAllBookingsInput{
		Page:     1,
		PageSize: 20,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 0).
		Return(testBookings, totalCount, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, testBookings, output.Bookings)
	assert.Equal(t, 1, output.Pagination.Page)
	assert.Equal(t, 20, output.Pagination.PageSize)
	assert.Equal(t, totalCount, output.Pagination.Total)
}

func TestGetAllBookings_Execute_Success_Page2(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	testBookings := []*bookings.Booking{
		createTestBooking("active"),
	}
	totalCount := 50

	input := GetAllBookingsInput{
		Page:     2,
		PageSize: 20,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 20).
		Return(testBookings, totalCount, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 2, output.Pagination.Page)
	assert.Equal(t, 20, output.Pagination.PageSize)
	assert.Equal(t, totalCount, output.Pagination.Total)
}

func TestGetAllBookings_Execute_Success_EmptyList(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	input := GetAllBookingsInput{
		Page:     1,
		PageSize: 20,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 0).
		Return([]*bookings.Booking{}, 0, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Empty(t, output.Bookings)
	assert.Equal(t, 0, output.Pagination.Total)
}

func TestGetAllBookings_Execute_StorageError(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	input := GetAllBookingsInput{
		Page:     1,
		PageSize: 20,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 0).
		Return(nil, 0, errors.New("database connection failed"))

	output, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "get bookings")
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestGetAllBookings_Execute_PageZero_DefaultsToOne(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	input := GetAllBookingsInput{
		Page:     0,
		PageSize: 20,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 0).
		Return([]*bookings.Booking{}, 0, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.Pagination.Page)
}

func TestGetAllBookings_Execute_PageNegative_DefaultsToOne(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	input := GetAllBookingsInput{
		Page:     -5,
		PageSize: 20,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 0).
		Return([]*bookings.Booking{}, 0, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.Pagination.Page)
}

func TestGetAllBookings_Execute_PageSizeZero_DefaultsToTwenty(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	input := GetAllBookingsInput{
		Page:     1,
		PageSize: 0,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 20, 0).
		Return([]*bookings.Booking{}, 0, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 20, output.Pagination.PageSize)
}

func TestGetAllBookings_Execute_PageSizeOver100_CapsTo100(t *testing.T) {
	bookingStorageMock := mocks.NewMockBookingStorage(t)
	uc := NewGetAllBookings(bookingStorageMock)

	input := GetAllBookingsInput{
		Page:     1,
		PageSize: 500,
	}

	bookingStorageMock.EXPECT().
		GetBookingsPaginated(mock.Anything, 100, 0).
		Return([]*bookings.Booking{}, 0, nil)

	output, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 100, output.Pagination.PageSize)
}

func TestGetAllBookings_Execute_OffsetCalculation(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantOffset int
	}{
		{"page 1, size 20", 1, 20, 0},  // (1-1)*20 = 0
		{"page 2, size 20", 2, 20, 20}, // (2-1)*20 = 20
		{"page 3, size 20", 3, 20, 40}, // (3-1)*20 = 40
		{"page 1, size 50", 1, 50, 0},  // (1-1)*50 = 0
		{"page 5, size 10", 5, 10, 40}, // (5-1)*10 = 40
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookingStorageMock := mocks.NewMockBookingStorage(t)
			uc := NewGetAllBookings(bookingStorageMock)

			bookingStorageMock.EXPECT().
				GetBookingsPaginated(mock.Anything, tt.pageSize, tt.wantOffset).
				Return([]*bookings.Booking{}, 0, nil)

			input := GetAllBookingsInput{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}

			output, err := uc.Execute(context.Background(), input)

			require.NoError(t, err)
			require.NotNil(t, output)
		})
	}
}

func TestGetAllBookings_Execute(t *testing.T) {
	tests := []struct {
		name            string
		input           GetAllBookingsInput
		storageBookings []*bookings.Booking
		storageTotal    int
		storageErr      error
		wantErr         bool
		wantPage        int
		wantPageSize    int
		wantTotal       int
	}{
		{
			name:            "success page 1",
			input:           GetAllBookingsInput{Page: 1, PageSize: 20},
			storageBookings: []*bookings.Booking{createTestBooking("active")},
			storageTotal:    50,
			storageErr:      nil,
			wantErr:         false,
			wantPage:        1,
			wantPageSize:    20,
			wantTotal:       50,
		},
		{
			name:            "success page 2",
			input:           GetAllBookingsInput{Page: 2, PageSize: 20},
			storageBookings: []*bookings.Booking{createTestBooking("active")},
			storageTotal:    50,
			storageErr:      nil,
			wantErr:         false,
			wantPage:        2,
			wantPageSize:    20,
			wantTotal:       50,
		},
		{
			name:            "page 0 defaults to 1",
			input:           GetAllBookingsInput{Page: 0, PageSize: 20},
			storageBookings: []*bookings.Booking{},
			storageTotal:    0,
			storageErr:      nil,
			wantErr:         false,
			wantPage:        1,
			wantPageSize:    20,
			wantTotal:       0,
		},
		{
			name:            "page size 0 defaults to 20",
			input:           GetAllBookingsInput{Page: 1, PageSize: 0},
			storageBookings: []*bookings.Booking{},
			storageTotal:    0,
			storageErr:      nil,
			wantErr:         false,
			wantPage:        1,
			wantPageSize:    20,
			wantTotal:       0,
		},
		{
			name:            "page size > 100 caps to 100",
			input:           GetAllBookingsInput{Page: 1, PageSize: 500},
			storageBookings: []*bookings.Booking{},
			storageTotal:    0,
			storageErr:      nil,
			wantErr:         false,
			wantPage:        1,
			wantPageSize:    100,
			wantTotal:       0,
		},
		{
			name:            "storage error",
			input:           GetAllBookingsInput{Page: 1, PageSize: 20},
			storageBookings: nil,
			storageTotal:    0,
			storageErr:      errors.New("db error"),
			wantErr:         true,
			wantPage:        0,
			wantPageSize:    0,
			wantTotal:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookingStorageMock := mocks.NewMockBookingStorage(t)
			uc := NewGetAllBookings(bookingStorageMock)

			expectedPageSize := tt.input.PageSize
			if expectedPageSize <= 0 {
				expectedPageSize = 20
			}
			if expectedPageSize > 100 {
				expectedPageSize = 100
			}
			expectedPage := tt.input.Page
			if expectedPage <= 0 {
				expectedPage = 1
			}
			expectedOffset := (expectedPage - 1) * expectedPageSize

			bookingStorageMock.EXPECT().
				GetBookingsPaginated(mock.Anything, expectedPageSize, expectedOffset).
				Return(tt.storageBookings, tt.storageTotal, tt.storageErr)

			output, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				require.NotNil(t, output)
				assert.Equal(t, tt.wantPage, output.Pagination.Page)
				assert.Equal(t, tt.wantPageSize, output.Pagination.PageSize)
				assert.Equal(t, tt.wantTotal, output.Pagination.Total)
			}
		})
	}
}
