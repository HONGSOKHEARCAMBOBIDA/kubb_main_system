package model

import "mysql/model/base"

type AcademicSection struct {
	base.ModelBase
	base.UUIDBase
	MajorID     int    `json:"major_id" gorm:"column:major_id"`
	ShiftID     int    `json:"shift_id" gorm:"column:shift_id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Type        int    `json:"type" gorm:"column:type"`
	Active      bool   `json:"active" gorm:"column:active"`
}

func (AcademicSection) TableName() string {
	return "academic_sections"
}
