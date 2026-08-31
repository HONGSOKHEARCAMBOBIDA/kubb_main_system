package response

import (
	"mysql/model"
	"mysql/model/base"
)

type AdmissionResponse struct {
	base.ModelBase
	base.UUIDBase
	StudentID                 int                  `json:"student_id"`
	StudentNameKh             string               `json:"student_name_kh"`
	StudentNameEn             string               `json:"student_name_en"`
	StudentGender             string               `json:"student_gender"`
	StudentGroupName          string               `json:"group_name" gorm:"column:group_name"`
	StudentDiscountType       model.DiscountType   `json:"discount_type"`
	StudentDiscountPercentage float64              `json:"discount_percentage" gorm:"column:discount_percentage"`
	StudentDiscountAmount     float64              `json:"discount_amount" gorm:"column:discount_amount"`
	TermID                    int                  `json:"term_id"`
	TermName                  string               `json:"term_name"`
	GenerationCode            string               `json:"generation_code"`
	GenerationName            string               `json:"generation_name"`
	AcademicCode              string               `json:"academic_code"`
	AcademicName              string               `json:"academic_name"`
	AcademicDegreeID          int                  `json:"academic_degree_id"`
	MajorCode                 string               `json:"major_code"`
	MajorName                 string               `json:"major_name"`
	ProgrammeID               int                  `json:"programme_id"`
	ProgrammeName             string               `json:"programme_name"`
	Name                      string               `json:"name" gorm:"column:name"`
	MonthlyFee                float64              `json:"monthly_fee" gorm:"column:monthly_fee"`
	QuarterlyFee              float64              `json:"quarterly_fee" gorm:"column:quarterly_fee"`
	SemesterlyFee             float64              `json:"semesterly_fee" gorm:"column:semesterly_fee"`
	YearlyFee                 float64              `json:"yearly_fee" gorm:"column:yearly_fee"`
	Date                      string               `json:"date" gorm:"column:date"`
	State                     model.AdmissionState `json:"state"`
	Description               string               `json:"description" gorm:"column:description"`
	ReferralSchool            string               `json:"referral_school" gorm:"column:referral_school"`
	Active                    bool                 `json:"active" gorm:"column:active"`
	EnrollmentResponse        []EnrollmentResponse `json:"enrollment" gorm:"-"`
}
