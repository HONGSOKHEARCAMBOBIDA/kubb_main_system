package model

import "mysql/model/base"

type Schoolarship struct {
	base.ModelBase
	base.UUIDBase
	Code               string       `json:"code" gorm:"column:code"`
	Name               string       `json:"name" gorm:"column:name"`
	DiscountType       DiscountType `gorm:"type:enum('percentage','amount');not null"`
	DiscountPercentage float64      `json:"	" gorm:"column:discount_percentage"`
	DiscountAmount     float64      `json:"discount_amount" gorm:"column:discount_amount"`
	Description        string       `json:"description" gorm:"column:description"`
	Active             bool         `json:"active" gorm:"column:active"`
}

func (Schoolarship) TableName() string {
	return "scholarships"
}
