package model

import "mysql/model/base"

type AcademicDegree struct {
	base.ModelBase
	base.UUIDBase
	AcademicID    int     `json:"academic_id" gorm:"column:academic_id"`
	MajorID       int     `json:"major_id" gorm:"column:major_id"`
	Name          string  `json:"name" gorm:"column:name"`
	MonthlyFee    float64 `json:"monthly_fee" gorm:"column:monthly_fee"`
	QuarterlyFee  float64 `json:"quarterly_fee" gorm:"column:quarterly_fee"`
	SemesterlyFee float64 `json:"semesterly_fee" gorm:"column:semesterly_fee"`
	YearlyFee     float64 `json:"yearly_fee" gorm:"column:yearly_fee"`
	Description   string  `json:"description" gorm:"column:description"`
	Active        bool    `json:"active" gorm:"column:active"`
}

func (AcademicDegree) TableName() string {
	return "academic_degrees"
}
