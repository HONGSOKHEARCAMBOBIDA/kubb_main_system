package model

import "mysql/model/base"

type StudentCategory struct {
	base.ModelBase
	Name string `json:"name" gorm:"column:name"`
}

func (StudentCategory) TableName() string {
	return "student_category"
}
