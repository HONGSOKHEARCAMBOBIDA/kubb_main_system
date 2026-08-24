package model

import "mysql/model/base"

type CourseRegistrations struct {
	base.ModelBase
	base.UUIDBase
	StudentTermID    int    `json:"student_term_id" gorm:"column:student_term_id"`
	ClassOfferingID  int    `json:"class_offering_id" gorm:"column:class_offering_id"`
	RegistrationDate string `json:"registration_date" gorm:"column:registration_date"`
	Status           string `json:"status" gorm:"column:status"`
}
