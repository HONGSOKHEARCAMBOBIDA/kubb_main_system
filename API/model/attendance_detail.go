package model

import (
	"mysql/model/base"
)

type AttendanceDetail struct {
	base.ModelBase
	base.UUIDBase
	AttendanceID          int              `json:"attendance_id" gorm:"column:attendance_id"`
	CourseRegistrationsID int              `json:"course_registration_id" gorm:"column:course_registration_id"`
	Status                AttendanceStatus `gorm:"type:enum('present','absent','late','excused');not null"`
	Note                  string           `json:"note" gorm:"column:note"`
}

func (AttendanceDetail) TableName() string {
	return "attendance_details"
}
