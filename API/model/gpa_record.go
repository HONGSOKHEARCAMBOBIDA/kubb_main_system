package model

import "mysql/model/base"

type GpaRecord struct {
	base.ModelBase
	base.UUIDBase
	StudentTermID int     `json:"student_term_id" gorm:"column:student_term_id"`
	TotalCredit   float64 `json:"total_credit" gorm:"column:total_credit"`
	SemesterGpa   float64 `json:"semester_gpa" gorm:"column:semester_gpa"`
	CumulativeGpa float64 `json:"cumulative_gpa" gorm:"column:cumulative_gpa"`
}

func (GpaRecord) TableName() string {
	return "gpa_records"
}
