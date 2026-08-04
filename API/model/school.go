package model

import "mysql/model/base"

type School struct {
	base.ModelBase
	Name         string `json:"name" gorm:"column:name"`
	Abbreviation string `json:"abbreviation" gorm:"column:abbreviation"`
	Director     string `json:"director" gorm:"column:director"`
	Phone        string `json:"phone" gorm:"column:phone"`
	Email        string `json:"email" gorm:"column:email"`
	Address      string `json:"address" gorm:"column:address"`
	Website      string `json:"website" gorm:"column:website"`
	Facebook     string `json:"facebook" gorm:"column:facebook"`
}

func (School) TableName() string {
	return "schools"
}
