package model

import "mysql/model/base"

type Schedule struct {
	base.ModelBase
	base.UUIDBase
	TeacherRateID  int            `json:"teacher_rate_id" gorm:"column:teacher_rate_id"`
	ScheduleDate   string         `json:"schedule_date" gorm:"column:schedule_date"`
	StartTime      string         `json:"start_time" gorm:"column:start_time"`
	EndTime        string         `json:"end_time" gorm:"column:end_time"`
	TotalTeachHour float64        `json:"total_teach_hours" gorm:"column:total_teach_hours"`
	Description    string         `json:"description" gorm:"column:description"`
	Active         bool           `json:"active" gorm:"column:active"`
	RoomID         int            `json:"room_id" gorm:"column:room_id"`
	Status         ScheduleStatus `gorm:"type:enum('active','cancelled','completed');not null"`
	VerifyBy       *int           `json:"verify_by" gorm:"column:verify_by"`
}
