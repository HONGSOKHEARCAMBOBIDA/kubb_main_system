package model

import "mysql/model/base"

type ClassCurriculumn struct {
	base.ModelBase
	base.UUIDBase
	Name    string `json:"name" gorm:"column:name"`
	MajorID int    `json:"major_id" gorm:"column:major_id"`
	TermID  int    `json:"term_id" gorm:"column:term_id"`
	Active  bool   `json:"active" gorm:"column:active"`
}

func (ClassCurriculumn) TableName() string {
	return "class_curriculums"
}
