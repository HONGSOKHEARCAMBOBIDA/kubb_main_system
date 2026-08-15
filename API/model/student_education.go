package model

import "mysql/model/base"

type StudentEducation struct {
	base.ModelBase
	base.UUIDBase
	StudentID       int    `json:"student_id" gorm:"column:student_id"`
	Level           string `json:"level" gorm:"column:level"`
	SchoolName      string `json:"school_name" gorm:"column:school_name"`
	VillageID       int    `json:"village_id" gorm:"column:village_id"`
	StartDate       string `json:"start_date" gorm:"column:start_date"`
	EndDate         string `json:"end_date" gorm:"column:end_date"`
	CertificateDate string `json:"cerificate_date" gorm:"column:cerificate_date"`
	Score           string `json:"score" gorm:"column:score"`
	Gpa             string `json:"gpa" gorm:"column:gpa"`
	Grade           string `json:"grade" gorm:"column:grade"`
}

func (StudentEducation) TableName() string {
	return "student_education"
}
