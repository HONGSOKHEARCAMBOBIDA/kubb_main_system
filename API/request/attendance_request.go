package request

import "mysql/model"

type AttendanceRequestCreate struct {
	ScheduleID                    int                             `json:"schedule_id" gorm:"column:schedule_id"`
	AttendanceDate                string                          `json:"attendance_date" gorm:"column:attendance_date"`
	Topic                         *string                         `json:"topic" gorm:"column:topic"`
	AttendanceDetailRequestCreate []AttendanceDetailRequestCreate `json:"detail"`
}

type AttendanceDetailRequestCreate struct {
	CourseRegistrationsID int                    `json:"course_registration_id" gorm:"column:course_registration_id"`
	Status                model.AttendanceStatus `json:"status"`
	Note                  *string                `json:"note" gorm:"column:note"`
}
