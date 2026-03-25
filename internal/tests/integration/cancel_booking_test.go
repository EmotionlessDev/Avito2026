package integration

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	bookingDto "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/dto"
	bookingStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/storage"
	bookingUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/usecases"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/tests/fixtures"
)

func (s *IntegrationSuite) TestCancelBooking_Success() {
	s.T().Log("Test: CancelBooking_Success")

	userFixture := fixtures.AdminUserFixture()
	userID := s.CreateUser(userFixture)

	roomFixture := fixtures.DefaultRoomFixture()
	roomID := s.CreateRoom(roomFixture)

	slotFixture := fixtures.DefaultSlotFixture()
	slotFixture.RoomID = roomID
	slotID := s.CreateSlot(slotFixture)

	bookingFixture := fixtures.DefaultBookingFixture()
	bookingFixture.SlotID = slotID
	bookingFixture.UserID = userID
	bookingFixture.Status = "active"
	bookingID := s.CreateBooking(bookingFixture)

	s.T().Logf("Created booking: %s", bookingID)

	bookingStorage := bookingStorage.NewStorage(s.db)
	cancelBookingUC := bookingUsecase.NewCancelBooking(bookingStorage)

	input := bookingDto.CancelBookingInput{
		BookingID: bookingID,
		UserID:    userID,
	}

	output, err := cancelBookingUC.Execute(s.ctx, input)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), output)
	require.NotNil(s.T(), output.Booking)
	assert.Equal(s.T(), bookingID, output.Booking.ID)
	assert.Equal(s.T(), "cancelled", output.Booking.Status)

	s.verifyBookingStatus(bookingID, "cancelled")

	s.T().Log("Test: CancelBooking_Success passed")
}

func (s *IntegrationSuite) TestCancelBooking_AlreadyCancelled() {
	s.T().Log("Test: CancelBooking_AlreadyCancelled")

	userFixture := fixtures.DefaultUserFixture()
	userID := s.CreateUser(userFixture)

	roomFixture := fixtures.DefaultRoomFixture()
	roomID := s.CreateRoom(roomFixture)

	slotFixture := fixtures.DefaultSlotFixture()
	slotFixture.RoomID = roomID
	slotID := s.CreateSlot(slotFixture)

	bookingFixture := fixtures.DefaultBookingFixture()
	bookingFixture.SlotID = slotID
	bookingFixture.UserID = userID
	bookingFixture.Status = "cancelled"
	bookingID := s.CreateBooking(bookingFixture)

	s.T().Logf("Created booking: %s", bookingID)

	bookingStorage := bookingStorage.NewStorage(s.db)
	cancelBookingUC := bookingUsecase.NewCancelBooking(bookingStorage)

	input := bookingDto.CancelBookingInput{
		BookingID: bookingID,
		UserID:    userID,
	}

	output, err := cancelBookingUC.Execute(s.ctx, input)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), output)
	require.NotNil(s.T(), output.Booking)
	assert.Equal(s.T(), bookingID, output.Booking.ID)
	assert.Equal(s.T(), "cancelled", output.Booking.Status)
	s.verifyBookingStatus(bookingID, "cancelled")

	s.T().Log("Test: CancelBooking_AlreadyCancelled passed")
}

func (s *IntegrationSuite) TestCancelBooking_DifferentUser() {
	s.T().Log("Test: CancelBooking_DifferentUser")

	userFixture := fixtures.AdminUserFixture()
	userID := s.CreateUser(userFixture)

	otherUserFixture := fixtures.DefaultUserFixture()
	otherUserID := s.CreateUser(otherUserFixture)

	roomFixture := fixtures.DefaultRoomFixture()
	roomID := s.CreateRoom(roomFixture)

	slotFixture := fixtures.DefaultSlotFixture()
	slotFixture.RoomID = roomID
	slotID := s.CreateSlot(slotFixture)

	bookingFixture := fixtures.DefaultBookingFixture()
	bookingFixture.SlotID = slotID
	bookingFixture.UserID = userID
	bookingFixture.Status = "active"
	bookingID := s.CreateBooking(bookingFixture)

	s.T().Logf("Created booking: %s", bookingID)

	bookingStorage := bookingStorage.NewStorage(s.db)
	cancelBookingUC := bookingUsecase.NewCancelBooking(bookingStorage)

	input := bookingDto.CancelBookingInput{
		BookingID: bookingID,
		UserID:    otherUserID,
	}

	output, err := cancelBookingUC.Execute(s.ctx, input)

	require.Error(s.T(), err)
	require.Nil(s.T(), output)

	s.verifyBookingStatus(bookingID, "active")

	s.T().Log("Test: CancelBooking_DifferentUser passed")
}

func (s *IntegrationSuite) TestCancelBooking_NonExistentBooking() {
	s.T().Log("Test: CancelBooking_NonExistentBooking")

	userFixture := fixtures.DefaultUserFixture()
	userID := s.CreateUser(userFixture)

	bookingStorage := bookingStorage.NewStorage(s.db)
	cancelBookingUC := bookingUsecase.NewCancelBooking(bookingStorage)

	input := bookingDto.CancelBookingInput{
		BookingID: uuid.New().String(),
		UserID:    userID,
	}

	output, err := cancelBookingUC.Execute(s.ctx, input)

	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), common.ErrBookingNotFound.Error())
	require.Nil(s.T(), output)

	s.T().Log("Test: CancelBooking_NonExistentBooking passed")
}

func (s *IntegrationSuite) verifyBookingStatus(bookingID, expectedStatus string) {
	s.T().Helper()

	query := `SELECT status FROM bookings WHERE id = $1`

	var actualStatus string
	err := s.db.QueryRow(query, bookingID).Scan(&actualStatus)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), expectedStatus, actualStatus)
}
