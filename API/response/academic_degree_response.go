package response

import "mysql/model/base"

type AcademicDegreeResponse struct {
	base.ModelBase
	base.UUIDBase
	AcademicID    int     `json:"academic_id" gorm:"column:academic_id"`
	AcademicCode  string  `json:"academic_code"`
	AcademicName  string  `json:"academic_name"`
	FacultyID     int     `json:"faculty_id"`
	DepartmentID  int     `json:"department_id"`
	MajorID       int     `json:"major_id" gorm:"column:major_id"`
	MajorCode     string  `json:"major_code"`
	MajorName     string  `json:"major_name"`
	ProgrammeID   int     `json:"programme_id"`
	ProgrammeName string  `json:"programme_name"`
	Name          string  `json:"name" gorm:"column:name"`
	MonthlyFee    float64 `json:"monthly_fee" gorm:"column:monthly_fee"`
	QuarterlyFee  float64 `json:"quarterly_fee" gorm:"column:quarterly_fee"`
	SemesterlyFee float64 `json:"semesterly_fee" gorm:"column:semesterly_fee"`
	YearlyFee     float64 `json:"yearly_fee" gorm:"column:yearly_fee"`
	Description   string  `json:"description" gorm:"column:description"`
	Active        bool    `json:"active" gorm:"column:active"`
}
