package model

import "mysql/model/base"

type Academic struct {
	base.ModelBase
	UUID        string  `json:"uuid" gorm:"column:uuid"`
	Code        string  `json:"code" gorm:"column:code"`
	Name        string  `json:"name" gorm:"column:name"`
	StartDate   string  `json:"start_date" gorm:"column:start_date"`
	EndDate     *string `json:"end_date" gorm:"column:end_date"`
	Description string  `json:"description" gorm:"column:description"`
	Active      bool    `json:"active" gorm:"column:active"`
}

func (Academic) TableName() string {
	return "academics"
}
