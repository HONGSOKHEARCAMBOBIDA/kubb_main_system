package response

import (
	"mysql/model"
	"mysql/model/base"
)

type EnrollmentResponse struct {
	base.ModelBase
	base.UUIDBase
	AdmissionID                    int                   `json:"admission_id" gorm:"column:admission_id"`
	SchoolarshipID                 int                   `json:"schoolarship_id"`
	SchoolarshipName               string                `json:"schoolarship_name"`
	SchoolarshipDiscountType       model.DiscountType    `json:"schoolarship_discount_type"`
	SchoolarshipDiscountPercentage float64               `json:"schoolarship_discount_percentage" gorm:"column:schoolarship_discount_percentage"`
	SchoolarshipDiscountAmount     float64               `json:"schoolarship_discount_amount" gorm:"column:schoolarship_discount_amount"`
	FeeInterval                    model.FeeInterval     `json:"fee_interval"`
	Description                    *string               `json:"description" gorm:"column:description"`
	StudentResponse                []StudentTermResponse `json:"student_term" gorm:"-"`
	FeeResponse                    []FeeResponse         `json:"fee_response" gorm:"-"`
	YearID                         int                   `json:"year_id"`
}
