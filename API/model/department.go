package model

import "mysql/model/base"

type Department struct {
	base.ModelBase
	base.UUIDBase
	FacultyID   int    `json:"faculty_id" gorm:"column:faculty_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Active      bool   `json:"active" gorm:"column:active"`
}

func (Department) TableName() string {
	return "departments"
}
