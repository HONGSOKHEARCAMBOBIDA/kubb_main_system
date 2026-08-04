package model

import "mysql/model/base"

type Campuse struct {
	base.ModelBase
	SchoolID        int    `json:"school_id" gorm:"column:school_id"`
	Code            string `json:"code" gorm:"column:code"`
	Prefix          string `json:"prefix" gorm:"column:prefix"`
	Name            string `json:"name" gorm:"column:name"`
	Director        string `json:"director" gorm:"column:director"`
	CapitalProvince string `json:"capital_province" gorm:"column:capital_province"`
	SchoolLevel     string `json:"school_level" gorm:"column:school_level"`
	Phone           string `json:"phone" gorm:"column:phone"`
	Email           string `json:"email" gorm:"column:email"`
	Website         string `json:"website" gorm:"column:website"`
	Facebook        string `json:"facebook" gorm:"column:facebook"`
	Address         string `json:"address" gorm:"column:address"`
	School          School `json:"school"`
}
