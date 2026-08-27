package model

import "mysql/model/base"

type Attendance struct {
	base.ModelBase
	base.UUIDBase
	ScheduleID     int    `json:"schedule_id" gorm:"column:schedule_id"`
	AttendanceDate string `json:"attendance_date" gorm:"column:attendance_date"`
	Topic          string `json:"topic" gorm:"column:topic"`
}

func (Attendance) TableName() string {
	return "attendances"
}
