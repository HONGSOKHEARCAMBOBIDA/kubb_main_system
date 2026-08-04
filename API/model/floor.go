package model

import "mysql/model/base"

type Floor struct {
	base.ModelBase
	BuildingID  int      `json:"building_id" gorm:"column:building_id"`
	Code        string   `json:"code" gorm:"column:code"`
	Name        string   `json:"name" gorm:"column:name"`
	Description string   `json:"description" gorm:"column:description"`
	Active      bool     `json:"active" gorm:"column:active"`
	Building    Building `json:"building"`
}

func (Floor) TableName() string {
	return "floors"
}
