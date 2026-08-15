package model

import "mysql/model/base"

type AcademicStream struct {
	base.ModelBase
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

func (AcademicStream) TableName() string {
	return "academic_streams"
}
