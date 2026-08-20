package response

import (
	"mysql/model"
	"mysql/model/base"
)

type EnrollmentResponse struct {
	base.ModelBase
	base.UUIDBase
	AdmissionID          int                   `json:"admission_id" gorm:"column:admission_id"`
	SchoolarshipResponse SchoolarshipResponse  `json:"schoolarship" gorm:"-"`
	FeeInterval          model.FeeInterval     `json:"fee_interval"`
	Description          *string               `json:"description" gorm:"column:description"`
	StudentResponse      []StudentTermResponse `json:"student_term" gorm:"-"`
}
