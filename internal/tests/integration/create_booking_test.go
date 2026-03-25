package integration

import (
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	bookingStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/storage"
	bookingUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/bookings/usecases"
	roomStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms/storage"
	roomUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms/usecases"
	scheduleStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/storage"
	scheduleUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/usecases"
	slotStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots/storage"
	slotUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/slots/usecases"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/tests/fixtures"
)

func (s *IntegrationSuite) TestCreateBooking_Success() {
	s.T().Log("Test: CreateBooking_Success")

	user := fixtures.DefaultUserFixture()
	userID := fixtures.CreateUser(s.T(), s.db, user)

	roomSt := roomStorage.NewStorage(s.db)
	scheduleSt := scheduleStorage.NewStorage(s.db)
	slotSt := slotStorage.NewStorage(s.db)
	bookingSt := bookingStorage.NewStorage(s.db)

	createRoomUC := roomUsecase.NewCreateRoom(roomSt)
	roomID, err := createRoomUC.Execute(s.ctx, "Test Room", "A room for testing bookings", 10)
	require.NoError(s.T(), err)

	createScheduleUC := scheduleUsecase.NewCreateSchedule(scheduleSt)
	_, err = createScheduleUC.Execute(s.ctx, scheduleUsecase.CreateScheduleInput{
		RoomID:     roomID,
		StartTime:  "09:00",
		EndTime:    "18:00",
		DaysOfWeek: []int{1, 2, 3, 4, 5},
	})
	require.NoError(s.T(), err)

	futureDate := time.Now().UTC().Add(7 * 24 * time.Hour)
	for futureDate.Weekday() == time.Saturday || futureDate.Weekday() == time.Sunday {
		futureDate = futureDate.Add(24 * time.Hour)
	}

	// get slots (and generate them if not exist)
	getSlotsUC := slotUsecase.NewGetSlots(scheduleSt, slotSt, roomSt)
	slots, err := getSlotsUC.Execute(s.ctx, slotUsecase.GetSlotsInput{
		RoomID: roomID,
		Date:   futureDate.Format("2006-01-02"),
	})

	require.NoError(s.T(), err)
	require.Greater(s.T(), len(slots), 0)

	slotToBook := slots[0]

	createBookingUC := bookingUsecase.NewCreateBooking(bookingSt, slotSt)

	input := bookingUsecase.CreateBookingInput{
		SlotID:               slotToBook.ID,
		UserID:               userID,
		CreateConferenceLink: true,
	}

	// Act
	booking, err := createBookingUC.Execute(s.ctx, input)

	// Assert
	require.NoError(s.T(), err)
	require.NotNil(s.T(), booking)
	assert.Equal(s.T(), slotToBook.ID, booking.SlotID)
	assert.Equal(s.T(), userID, booking.UserID)
	assert.Equal(s.T(), "active", booking.Status)
	assert.NotEmpty(s.T(), booking.ConferenceLink)

	s.T().Log("Test: CreateBooking_Success passed")
}

func (s *IntegrationSuite) TestCreateBooking_SlotAlreadyBooked() {
	s.T().Log("Test: CreateBooking_SlotAlreadyBooked")

	user := fixtures.DefaultUserFixture()
	userID := fixtures.CreateUser(s.T(), s.db, user)

	roomSt := roomStorage.NewStorage(s.db)
	scheduleSt := scheduleStorage.NewStorage(s.db)
	slotSt := slotStorage.NewStorage(s.db)
	bookingSt := bookingStorage.NewStorage(s.db)

	createRoomUC := roomUsecase.NewCreateRoom(roomSt)
	roomID, err := createRoomUC.Execute(s.ctx, "Test Room", "A room for testing bookings", 10)
	require.NoError(s.T(), err)

	createScheduleUC := scheduleUsecase.NewCreateSchedule(scheduleSt)
	_, err = createScheduleUC.Execute(s.ctx, scheduleUsecase.CreateScheduleInput{
		RoomID:     roomID,
		StartTime:  "09:00",
		EndTime:    "18:00",
		DaysOfWeek: []int{1, 2, 3, 4, 5},
	})
	require.NoError(s.T(), err)

	futureDate := time.Now().UTC().Add(7 * 24 * time.Hour)
	for futureDate.Weekday() == time.Saturday || futureDate.Weekday() == time.Sunday {
		futureDate = futureDate.Add(24 * time.Hour)
	}

	getSlotsUC := slotUsecase.NewGetSlots(scheduleSt, slotSt, roomSt)
	slots, err := getSlotsUC.Execute(s.ctx, slotUsecase.GetSlotsInput{
		RoomID: roomID,
		Date:   futureDate.Format("2006-01-02"),
	})

	require.NoError(s.T(), err)
	require.Greater(s.T(), len(slots), 0)

	slotToBook := slots[0]

	createBookingUC := bookingUsecase.NewCreateBooking(bookingSt, slotSt)

	input := bookingUsecase.CreateBookingInput{
		SlotID:               slotToBook.ID,
		UserID:               userID,
		CreateConferenceLink: true,
	}

	// Act
	_, err = createBookingUC.Execute(s.ctx, input)

	require.NoError(s.T(), err)

	// Try to book the same slot again
	_, err = createBookingUC.Execute(s.ctx, input)

	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), common.ErrSlotAlreadyBooked.Error())
	s.T().Log("Test: CreateBooking_SlotAlreadyBooked passed")
}
