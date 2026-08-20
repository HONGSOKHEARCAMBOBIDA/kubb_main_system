package response

import "mysql/model/base"

type StudentTermResponse struct {
	base.ModelBase
	base.UUIDBase
	EnrollmentID      int               `json:"enrollment_id" gorm:"column:enrollment_id"`
	SemesternResponse SemesternResponse `json:"semester" gorm:"-"`
	FeeResponse       []FeeResponse     `json:"fee" gorm:"-"`
	StudyYearID       int               `json:"study_year_id" gorm:"column:study_year_id"`
	Active            bool              `json:"active" gorm:"column:active"`
	Status            string            `json:"status" gorm:"column:status"`
}
