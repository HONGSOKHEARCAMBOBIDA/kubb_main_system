package model

import "mysql/model/base"

type Enrollment struct {
	base.ModelBase
	base.UUIDBase
	AdmissionID    int         `json:"admission_id" gorm:"column:admission_id"`
	SchoolarshipID int         `json:"scholarship_id" gorm:"column:scholarship_id"`
	SectionID      *int        `json:"section_id" gorm:"column:section_id"`
	FeeInterval    FeeInterval `gorm:"type:enum('monthly_fee','quarterly_fee','semesterly_fee','yearly_fee');not null"`
	Description    *string     `json:"description" gorm:"column:description"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}
