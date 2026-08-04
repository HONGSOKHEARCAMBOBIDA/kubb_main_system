package model

import "mysql/model/base"

type Building struct {
	base.ModelBase
	Code        string  `json:"code" gorm:"column:code"`
	CampuseID   int     `json:"campus_id" gorm:"column:campus_id"`
	Name        string  `json:"name" gorm:"column:name"`
	Address     string  `json:"address" gorm:"column:address"`
	Description string  `json:"description" gorm:"column:description"`
	Campuse     Campuse `json:"campuse"`
}

func (Building) TableName() string {
	return "buildings"
}
