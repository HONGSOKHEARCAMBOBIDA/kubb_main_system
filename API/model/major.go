package model

import "mysql/model/base"

type Major struct {
	base.ModelBase
	base.UUIDBase
	DepartmentID     int              `json:"department_id" gorm:"column:department_id"`
	Code             string           `json:"code" gorm:"column:code"`
	Name             string           `json:"name" gorm:"column:name"`
	DurationPeriod   int              `json:"duration_period" gorm:"column:duration_period"`
	DurationInterval DurationInterval `gorm:"type:enum('year','month','week','day');not null;default:'year'"`
	Description      string           `gorm:"type:longtext;column:description"`
	Active           bool             `gorm:"not null;default:true"`
}
