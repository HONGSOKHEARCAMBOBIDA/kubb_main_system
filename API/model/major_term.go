package model

import "mysql/model/base"

type MajorTerm struct {
	base.ModelBase
	base.UUIDBase
	MajorID int  `json:"major_id" gorm:"column:major_id"`
	TermID  int  `json:"term_id" gorm:"column:term_id"`
	Active  bool `json:"active" gorm:"column:active"`
}

func (MajorTerm) TableName() string {
	return "major_terms"
}
