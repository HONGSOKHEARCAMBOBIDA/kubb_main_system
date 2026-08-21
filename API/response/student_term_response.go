package response

import "mysql/model/base"

type StudentTermResponse struct {
	base.ModelBase
	base.UUIDBase
	EnrollmentID int    `json:"enrollment_id" gorm:"column:enrollment_id"`
	SemesterID   int    `json:"semester_id"`
	SemesterCode string `json:"semester_code"`
	SemesterName string `json:"semester_name"`
	AcademicID   int    `json:"academic_id" gorm:"column:academic_id"`
	AcademicName string `json:"academic_name"`
	//	FeeResponse  []FeeResponse `json:"fee" gorm:"-"`
	StudyYearID int    `json:"study_year_id" gorm:"column:study_year_id"`
	Active      bool   `json:"active" gorm:"column:active"`
	Status      string `json:"status" gorm:"column:status"`
}
