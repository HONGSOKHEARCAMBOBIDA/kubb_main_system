package model

import "mysql/model/base"

type TeacherRate struct {
	base.ModelBase
	base.UUIDBase
	TeacherID       int     `json:"teacher_id" gorm:"column:teacher_id"`
	ClassOfferingID int     `json:"class_offer_id" gorm:"column:class_offer_id"`
	HourlyRate      float64 `json:"hourly_rate" gorm:"column:hourly_rate"`
	EffectiveDate   string  `json:"effective_date" gorm:"column:effective_date"`
	EndDate         string  `json:"end_date" gorm:"column:end_date"`
	Active          bool    `json:"active" gorm:"column:active"`
	CreateBy        int     `json:"created_by" gorm:"column:created_by"`
}
