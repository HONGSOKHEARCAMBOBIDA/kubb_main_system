package request

import "mysql/model"

type AdmissionRequestCreate struct {
	TermID           int                  `json:"term_id" gorm:"column:term_id"`
	AcademicDegreeID int                  `json:"academic_degree_id" gorm:"column:academic_degree_id"`
	Date             string               `json:"date" gorm:"column:date"`
	AdmissionState   model.AdmissionState `json:"state"`
	Description      string               `json:"description" gorm:"column:description"`
	ReferralSchool   string               `json:"referral_school" gorm:"column:referral_school"`
}
