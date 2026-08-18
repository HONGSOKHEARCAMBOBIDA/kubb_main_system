package request

import "mysql/model"

type EnrollmentRequestCreate struct {
	SchoolarshipID int               `json:"scholarship_id" gorm:"column:scholarship_id"`
	FeeInterval    model.FeeInterval `json:"fee_interval"`
}
