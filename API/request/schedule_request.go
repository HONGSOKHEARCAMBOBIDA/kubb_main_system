package request

import "mysql/model"

type ScheduleRequest struct {
	TeacherRateID  int     `json:"teacher_rate_id" gorm:"column:teacher_rate_id"`
	ScheduleDate   string  `json:"schedule_date" gorm:"column:schedule_date"`
	StartTime      string  `json:"start_time" gorm:"column:start_time"`
	EndTime        string  `json:"end_time" gorm:"column:end_time"`
	TotalTeachHour float64 `json:"total_teach_hours" gorm:"column:total_teach_hours"`
	Description    string  `json:"description" gorm:"column:description"`
	Active         bool    `json:"active" gorm:"column:active"`
	RoomID         int     `json:"room_id" gorm:"column:room_id"`
}

type ScheduleRequestUpdate struct {
	TotalTeachHour float64              `json:"total_teach_hours" gorm:"column:total_teach_hours"`
	Status         model.ScheduleStatus `json:"status"`
}
