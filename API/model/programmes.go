package model

import "mysql/model/base"

type Programmes struct {
	base.ModelBase
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Active      bool   `json:"active" gorm:"column:active"`
}
