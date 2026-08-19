package model

import "mysql/model/base"

type Fee struct {
	base.ModelBase
	base.UUIDBase
	EnrollmentID int     `json:"enrollment_id" gorm:"column:enrollment_id"`
	Date         string  `json:"date" gorm:"column:date"`
	Amount       float64 `json:"amount" gorm:"column:amount"`
	Discount     float64 `json:"discount" gorm:"column:discount"`
	Total        float64 `json:"total" gorm:"column:total"`
	Active       bool    `json:"active" gorm:"column:active"`
}

func (Fee) TableName() string {
	return "fees"
}
