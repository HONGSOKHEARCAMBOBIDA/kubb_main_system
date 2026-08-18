package model

import "mysql/model/base"

type Admission struct {
	base.ModelBase
	base.UUIDBase
	StudentID        int            `json:"student_id" gorm:"column:student_id"`
	TermID           int            `json:"term_id" gorm:"column:term_id"`
	AcademicDegreeID int            `json:"academic_degree_id" gorm:"column:academic_degree_id"`
	Date             string         `json:"date" gorm:"column:date"`
	AdmissionState   AdmissionState `gorm:"type:enum('created','submitted','approved','rejected','cancelled');not null"`
	Description      string         `json:"description" gorm:"column:description"`
	ReferralSchool   string         `json:"referral_school" gorm:"column:referral_school"`
	Active           bool           `json:"active" gorm:"column:active"`
}
