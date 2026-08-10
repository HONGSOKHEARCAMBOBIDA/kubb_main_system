package response

import (
	"mysql/model"
	"mysql/model/base"
)

type MajorResponse struct {
	base.ModelBase
	base.UUIDBase
	Name             string                 `json:"name"`
	Code             string                 `json:"code"`
	DurationPeriod   int                    `json:"duration_period" gorm:"column:duration_period"`
	DurationInterval model.DurationInterval `json:"duration_interval"`
	Description      string                 `json:"description"`
	Active           bool                   `json:"active" gorm:"column:active"`
	DepartmentID     int                    `json:"department_id"`
	DepartmentName   string                 `json:"department_name"`
	DepartmentCode   string                 `json:"department_code"`
	FacultyID        int                    `json:"faculty_id"`
	FacultyName      string                 `json:"faculty_name"`
	FacultyCode      string                 `json:"faculty_code"`
	ProgrammeID      int                    `json:"programme_id"`
	ProgrammName     string                 `json:"programme_name" gorm:"column:programme_name"`
}

type MajorResponseByDepartment struct {
	base.ModelBase
	Name string `json:"name" gorm:"column:name"`
}
