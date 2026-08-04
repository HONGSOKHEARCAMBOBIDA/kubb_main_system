package model

import "mysql/model/base"

type SchoolOffice struct {
	base.ModelBase
	FloorID     int    `json:"floor_id" gorm:"column:floor_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Address     string `json:"address" gorm:"column:address"`
	Description string `json:"description" gorm:"column:description"`
	Active      bool   `json:"active" gorm:"column:active"`
	Floor       Floor  `json:"floor"`
}

func (SchoolOffice) TableName() string {
	return "school_offices"
}
