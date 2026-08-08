package response

import "mysql/model/base"

type FacultyResponse struct {
	base.ModelBase
	ProgrammeID   int    `json:"programme_id" gorm:"column:programme_id"`
	ProgrammeName string `json:"programme_name"`
	Code          string `json:"code" gorm:"column:code"`
	Name          string `json:"name" gorm:"column:name"`
	Description   string `json:"description" gorm:"column:description"`
	Active        bool   `json:"active" gorm:"column:active"`
}

type FacultyResponseByProgrammes struct {
	base.ModelBase
	Name string `json:"name" gorm:"column:name"`
}
