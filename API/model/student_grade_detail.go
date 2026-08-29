package model

import "mysql/model/base"

type StudentGradeDetail struct {
	base.ModelBase
	base.UUIDBase
	StudentGradeID   int      `json:"student_grade_id" gorm:"column:student_grade_id"`
	GradeComponentID int      `json:"grade_component_id" gorm:"column:grade_component_id"`
	Score            *float64 `json:"score" gorm:"column:score"`
}

func (StudentGradeDetail) TableName() string {
	return "student_grade_details"
}
