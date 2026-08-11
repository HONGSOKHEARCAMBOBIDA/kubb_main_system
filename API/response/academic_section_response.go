package response

import "mysql/model/base"

type AcademicSectionResponse struct {
	base.ModelBase
	base.UUIDBase
	MajorID       int    `json:"major_id" gorm:"column:major_id"`
	MajorCode     string `json:"major_code"`
	MajorName     string `json:"major_name"`
	ProgrammeID   int    `json:"programme_id"`
	ProgrammeName string `json:"programme_name"`
	ShiftID       int    `json:"shift_id" gorm:"column:shift_id"`
	ShiftName     string `json:"shift_name"`
	AcademicID    int    `json:"academic_id"`
	AcademicName  string `json:"academic_name"`
	Name          string `json:"name" gorm:"column:name"`
	Description   string `json:"description" gorm:"column:description"`
	Active        bool   `json:"active" gorm:"column:active"`
}
