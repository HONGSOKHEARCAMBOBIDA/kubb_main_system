package response

import (
	"mysql/model"
	"mysql/model/base"
)

type ClasCurriculumnResponse struct {
	base.ModelBase
	base.UUIDBase
	Name                          string                          `json:"name" gorm:"column:name"`
	Active                        bool                            `json:"active" gorm:"column:active"`
	MajorID                       int                             `json:"major_id"`
	MajorName                     string                          `json:"major_name"`
	MajorCode                     string                          `json:"major_code"`
	MajorDurationPeriod           int                             `json:"major_duration_period" gorm:"column:major_duration_period"`
	MajorDurationInterval         model.DurationInterval          `json:"major_duration_interval"`
	DepartmentID                  int                             `json:"department_id"`
	DepartmentName                string                          `json:"department_name"`
	DepartmentCode                string                          `json:"department_code"`
	FacultyName                   string                          `json:"faculty_name"`
	FacultyID                     int                             `json:"faculty_id"`
	FacultyCode                   string                          `json:"faculty_code"`
	ProgrammeID                   int                             `json:"programme_id"`
	ProgrammName                  string                          `json:"programme_name" gorm:"column:programme_name"`
	TermID                        int                             `json:"term_id"`
	TermCode                      string                          `json:"term_code"`
	TermName                      string                          `json:"term_name"`
	GenerationID                  int                             `json:"generation_id" gorm:"column:generation_id"`
	GenerationCode                string                          `json:"generation_code"`
	GenerationName                string                          `json:"generation_name"`
	AcademicID                    int                             `json:"academic_id"`
	AcademicCode                  string                          `json:"academic_code"`
	AcademicName                  string                          `json:"academic_name"`
	ClasCurriculumnDetailResponse []ClasCurriculumnDetailResponse `json:"class_curriculum_detais" gorm:"-"`
}
