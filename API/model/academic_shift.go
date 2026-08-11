package model

import "mysql/model/base"

type AcademicShift struct {
	base.ModelBase
	base.UUIDBase
	AcademicID  int    `json:"academic_id" gorm:"column:academic_id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Active      bool   `json:"active" gorm:"column:active"`
}

func (AcademicShift) TableName() string {
	return "academic_shifts"
}
