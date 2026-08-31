package request

import "mysql/model"

type EnrollmentRequestCreate struct {
	SchoolarshipID int               `json:"scholarship_id" gorm:"column:scholarship_id"`
	FeeInterval    model.FeeInterval `json:"fee_interval"`
}

type EnrollmentRequestCreateV2 struct {
	AdmissionID              int                       `json:"admision_id"`
	SchoolarshipID           int                       `json:"scholarship_id" gorm:"column:scholarship_id"`
	FeeInterval              model.FeeInterval         `json:"fee_interval"`
	StudentTermRequestCreate *StudentTermRequestCreate `json:"student_term"`
}

type EnrollmentRequestUpdate struct {
	Description string `json:"description"`
}
