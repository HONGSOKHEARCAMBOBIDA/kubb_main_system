package model

import "mysql/model/base"

type Term struct {
	base.ModelBase
	base.UUIDBase
	GenerationID int    `json:"generation_id" gorm:"column:generation_id"`
	Code         string `json:"code" gorm:"column:code"`
	Name         string `json:"name" gorm:"column:name"`
	Index        int    `json:"index" gorm:"column:index"`
	StartDate    string `json:"start_date" gorm:"column:start_date"`
	EndDate      string `json:"end_date" gorm:"column:end_date"`
	Description  string `json:"description" gorm:"column:description"`
	Active       bool   `json:"active" gorm:"column:active"`
}
