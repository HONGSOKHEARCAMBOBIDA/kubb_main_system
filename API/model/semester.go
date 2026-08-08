package model

import "mysql/model/base"

type Semester struct {
	base.ModelBase
	base.UUIDBase
	AcademicID int      `json:"academic_id" gorm:"column:academic_id"`
	Code       string   `json:"code" gorm:"column:code"`
	Name       string   `json:"name" gorm:"column:name"`
	Index      int      `json:"index" gorm:"column:index"`
	StartDate  string   `json:"start_date" gorm:"column:start_date"`
	EndDate    string   `json:"end_date" gorm:"column:end_date"`
	Active     bool     `json:"active" gorm:"column:active"`
	Academic   Academic `json:"academic"`
}

func (Semester) TableName() string {
	return "semesters"
}
