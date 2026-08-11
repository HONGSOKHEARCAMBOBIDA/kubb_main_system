package response

import "mysql/model/base"

type AcademicShiftResponse struct {
	base.ModelBase
	base.UUIDBase
	AcademicID   int    `json:"academic_id" gorm:"column:academic_id"`
	AcademicCode string `json:"academic_code"`
	AcademicName string `json:"academic_name"`
	Name         string `json:"name" gorm:"column:name"`
	Description  string `json:"description" gorm:"column:description"`
	Active       bool   `json:"active" gorm:"column:active"`
}

type AcademicShiftResponseByAcademic struct {
	base.ModelBase
	Name string `json:"name" gorm:"column:name"`
}
