package dto

type CreateScheduleInput struct {
	RoomID     string
	StartTime  string // "HH:MM"
	EndTime    string // "HH:MM"
	DaysOfWeek []int
}
