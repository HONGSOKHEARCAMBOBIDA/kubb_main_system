package model

type ScheduleStatus string

const (
	ScheduleStatusActive    ScheduleStatus = "active"
	ScheduleStatusCancell   ScheduleStatus = "cancelled"
	ScheduleStatusCompleted ScheduleStatus = "completed"
)
