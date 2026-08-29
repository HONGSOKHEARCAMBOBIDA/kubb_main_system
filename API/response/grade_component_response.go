package response

import "mysql/model/base"

type GradeComponentResponse struct {
	base.ModelBase
	base.UUIDBase
	SubjectID        int     `json:"subject_id" gorm:"column:subject_id"`
	Name             string  `json:"name" gorm:"column:name"`
	WeightPercentage float64 `json:"weight_percentage" gorm:"column:weight_percentage"`
	Active           bool    `json:"active" gorm:"column:active"`
}
