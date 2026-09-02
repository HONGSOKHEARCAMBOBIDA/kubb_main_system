package model

import "mysql/model/base"

type SubjectGroup struct {
	base.ModelBase
	base.UUIDBase
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}
