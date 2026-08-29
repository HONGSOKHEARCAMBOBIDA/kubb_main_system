package request

type GradeComponentRequestCreate struct {
	SubjectID             int                     `json:"subject_id" gorm:"column:subject_id"`
	GradeComponentRequest []GradeComponentRequest `json:"grade"`
}

type GradeComponentRequest struct {
	Name             string  `json:"name" gorm:"column:name"`
	WeightPercentage float64 `json:"weight_percentage" gorm:"column:weight_percentage"`
}
