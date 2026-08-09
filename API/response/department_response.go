package response

import "mysql/model/base"

type DepartmentResponse struct {
	base.UUIDBase
	base.ModelBase
	ProgrammeID   int    `json:"programme_id"`
	ProgrammeName string `json:"programme_name"`
	FacultyID     int    `json:"faculty_id"`
	FacultyCode   string `json:"faculty_code"`
	FacultyName   string `json:"faculty_name"`
	Code          string `json:"code" gorm:"column:code"`
	Name          string `json:"name" gorm:"column:name"`
	Description   string `json:"description" gorm:"column:description"`
	Active        bool   `json:"active" gorm:"column:active"`
}

type DepartmentResponseByFaculty struct {
	base.ModelBase
	Name string `json:"name" gorm:"column:name"`
}
