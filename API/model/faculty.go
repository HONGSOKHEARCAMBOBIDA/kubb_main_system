package model

import "mysql/model/base"

type Faculty struct {
	base.ModelBase
	base.UUIDBase
	ProgrammeID int    `json:"programme_id" gorm:"column:programme_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Active      bool   `json:"active" gorm:"column:active"`
}

func (Faculty) TableName() string {
	return "faculties"
}
