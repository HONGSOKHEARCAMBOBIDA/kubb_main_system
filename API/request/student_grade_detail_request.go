package request

type StudentGradeDetailRequestCreate struct {
	GradeComponentID int     `json:"grade_component_id" gorm:"column:grade_component_id"`
	Score            float64 `json:"score" gorm:"column:score"`
}
