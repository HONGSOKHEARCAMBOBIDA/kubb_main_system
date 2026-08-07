package response

import "mysql/model/base"

type TermResponse struct {
	base.ModelBase
	base.UUIDBase
	GenerationID   int    `json:"generation_id" gorm:"column:generation_id"`
	GenerationCode string `json:"generation_code"`
	GenerationName string `json:"generation_name"`
	AcademicID     int    `json:"academic_id"`
	AcademicCode   string `json:"academic_code"`
	AcademicName   string `json:"academic_name"`
	Code           string `json:"code" gorm:"column:code"`
	Name           string `json:"name" gorm:"column:name"`
	Index          int    `json:"index" gorm:"column:index"`
	StartDate      string `json:"start_date" gorm:"column:start_date"`
	EndDate        string `json:"end_date" gorm:"column:end_date"`
	Description    string `json:"description" gorm:"column:description"`
	Active         bool   `json:"active" gorm:"column:active"`
}
