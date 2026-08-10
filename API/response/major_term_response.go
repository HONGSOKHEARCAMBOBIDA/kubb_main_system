package response

import (
	"mysql/model"
	"mysql/model/base"
)

type MajorTermResponse struct {
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
	TermID           int                    `json:"term_id"`
	TermCode         string                 `json:"term_code"`
	TermName         string                 `json:"term_name"`
	GenerationID     int                    `json:"generation_id"`
	GenerationCode   string                 `json:"generation_code"`
	GenerationName   string                 `json:"generation_name"`
	AcademicID       int                    `json:"academic_id"`
	AcademicCode     string                 `json:"academic_code"`
	AcademicName     string                 `json:"academic_name"`
}
