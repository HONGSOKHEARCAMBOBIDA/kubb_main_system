package response

import "mysql/model/base"

type StudentGradeResponse struct {
	base.ModelBase
	base.UUIDBase
	CourseRegistrationID       int                          `json:"course_registration_id" gorm:"column:course_registration_id"`
	TotalScore                 *float64                     `json:"total_score" gorm:"column:total_score"`
	LetterGrade                *string                      `json:"letter_grade" gorm:"column:letter_grade"`
	GradePoint                 *float64                     `json:"grade_point" gorm:"column:grade_point"`
	Status                     string                       `json:"status"`
	StudentGradeDetailResponse []StudentGradeDetailResponse `json:"detail" gorm:"-"`
}

type StudentGradeDetailResponse struct {
	base.ModelBase
	base.UUIDBase
	StudentGradeID     int      `json:"student_grade_id" gorm:"column:student_grade_id"`
	GradeComponentID   int      `json:"grade_component_id" gorm:"column:grade_component_id"`
	GradeComponentName string   `json:"grade_component_name"`
	Score              *float64 `json:"score" gorm:"column:score"`
}
