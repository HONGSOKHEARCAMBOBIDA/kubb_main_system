package model

import "mysql/model/base"

type Fee struct {
	base.ModelBase
	base.UUIDBase
	StudentTermID int     `json:"student_term_id" gorm:"column:student_term_id"`
	Date          string  `json:"date" gorm:"column:date"`
	Amount        float64 `json:"amount" gorm:"column:amount"`
	Discount      float64 `json:"discount" gorm:"column:discount"`
	Total         float64 `json:"total" gorm:"column:total"`
	Active        bool    `json:"active" gorm:"column:discount"`
}
