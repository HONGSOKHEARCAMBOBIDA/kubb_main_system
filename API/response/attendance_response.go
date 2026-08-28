package response

import (
	"mysql/model"
	"mysql/model/base"
)

type AttendanceResponse struct {
	base.ModelBase
	base.UUIDBase
	AttendanceDate           string                     `json:"attendance_date" gorm:"column:attendance_date"`
	Topic                    string                     `json:"topic" gorm:"column:topic"`
	AttendanceDetailResponse []AttendanceDetailResponse `json:"detail" gorm:"-"`
}

type AttendanceDetailResponse struct {
	base.ModelBase
	base.UUIDBase
	AttendanceID         int                    `json:"attendance_id" gorm:"column:attendance_id"`
	CourseRegistrationID int                    `json:"course_registration_id"`
	NameKh               string                 `json:"name_kh"`
	NameEn               string                 `json:"name_en"`
	Dob                  string                 `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender               string                 `json:"gender"`
	Phone                string                 `json:"phone"`
	Status               model.AttendanceStatus `json:"status"`
	Note                 string                 `json:"note"`
}
