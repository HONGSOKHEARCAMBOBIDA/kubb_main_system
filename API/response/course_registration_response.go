package response

import "mysql/model/base"

type CourseRegistrationResponse struct {
	base.ModelBase
	base.UUIDBase
	NameKh string `json:"name_kh"`
	NameEn string `json:"name_en"`
	Dob    string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
}
