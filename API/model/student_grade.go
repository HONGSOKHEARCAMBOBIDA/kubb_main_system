package model

import "mysql/model/base"

type StudentGrade struct {
	base.ModelBase
	base.UUIDBase
	CourseRegistrationID int      `json:"course_registration_id" gorm:"column:course_registration_id"`
	TotalScore           *float64 `json:"total_score" gorm:"column:total_score"`
	LetterGrade          *string  `json:"letter_grade" gorm:"column:letter_grade"`
	GradePoint           *float64 `json:"grade_point" gorm:"column:grade_point"`
	Status               string   `json:"status"`
}

func (StudentGrade) TableName() string {
	return "student_grades"
}
