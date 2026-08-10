package response

import "mysql/model/base"

type SubjectResponse struct {
	base.ModelBase
	base.UUIDBase
	Code           string  `json:"code" gorm:"column:code"`
	Name           string  `json:"name" gorm:"column:name"`
	Credit         float64 `json:"credit" gorm:"column:credit"`
	PassingScore   float64 `json:"passing_score" gorm:"column:passing_score"`
	Description    string  `json:"description" gorm:"column:description"`
	Active         bool    `json:"active" gorm:"column:active"`
	MajorID        int     `json:"major_id" gorm:"column:major_id"`
	MajorCode      string  `json:"major_code"`
	MajorName      string  `json:"major_name"`
	DepartmentID   int     `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	DepartmentCode string  `json:"department_code"`
	FacultyID      int     `json:"faculty_id"`
	FacultyName    string  `json:"faculty_name"`
	FacultyCode    string  `json:"faculty_code"`
	ProgrammeID    int     `json:"programme_id"`
	ProgrammName   string  `json:"programme_name" gorm:"column:programme_name"`
}

type SubjectResponseByMajor struct {
	base.ModelBase
	Name string `json:"name" gorm:"column:name"`
}
