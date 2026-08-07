package response

import "mysql/model/base"

type GenerationResponse struct {
	base.ModelBase
	base.UUIDBase
	AcademicID   int     `json:"academic_id" gorm:"column:academic_id"`
	AcademicName string  `json:"academic_name"`
	Code         string  `json:"code" gorm:"column:code"`
	Name         string  `json:"name" gorm:"column:name"`
	Index        int     `json:"index" gorm:"column:index"`
	StartDate    string  `json:"start_date" gorm:"column:start_date"`
	EndDate      *string `json:"end_date" gorm:"column:end_date"`
	Description  string  `json:"description" gorm:"column:description"`
	Active       bool    `json:"active" gorm:"column:active"`
}

type GenerationResponseByAcademic struct {
	base.ModelBase
	Name string `json:"name"`
}
