package response

import "mysql/model/base"

type GpaRecordResponse struct {
	base.ModelBase
	base.UUIDBase
	StudentTermID int     `json:"student_term_id" gorm:"column:student_term_id"`
	TotalCredit   float64 `json:"total_credit" gorm:"column:total_credit"`
	SemesterGpa   float64 `json:"semester_gpa" gorm:"column:semester_gpa"`
	CumulativeGpa float64 `json:"cumulative_gpa" gorm:"column:cumulative_gpa"`
}
