package model

import "mysql/model/base"

type SchoolRoom struct {
	base.ModelBase
	FloorID     int    `json:"floor_id" gorm:"column:floor_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Capacity    int    `json:"capacity" gorm:"column:capacity"`
	Description string `json:"description" gorm:"column:description"`
	Active      bool   `json:"active" gorm:"column:active"`
	Floor       Floor  `json:"floor"`
}

func (SchoolRoom) TableName() string {
	return "school_rooms"
}
