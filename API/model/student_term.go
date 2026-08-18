package model

import "mysql/model/base"

type StudentTerm struct {
	base.ModelBase
	base.UUIDBase
	EnrollmentID int    `json:"enrollment_id" gorm:"column:enrollment_id"`
	SemesterID   int    `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID  int    `json:"study_year_id" gorm:"column:study_year_id"`
	Active       bool   `json:"active" gorm:"column:active"`
	Status       string `json:"status" gorm:"column:status"`
}

func (StudentTerm) TableName() string {
	return "student_terms"
}
