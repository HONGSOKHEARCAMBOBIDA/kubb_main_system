package model

import "mysql/model/base"

type ClassCurriculumnDetail struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnID int    `json:"class_curriculum_id" gorm:"column:class_curriculum_id"`
	SemesterID         int    `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID        int    `json:"study_year_id" gorm:"column:study_year_id"`
	AcademicShiftID    int    `json:"academic_shift_id" gorm:"column:academic_shift_id"`
	StartDate          string `json:"start_date" gorm:"column:start_date"`
	MidtermDate        string `json:"midterm_date" gorm:"column:midterm_date"`
	FinalDate          string `json:"final_date" gorm:"column:final_date"`
	TotalStudent       *int   `json:"total_student" gorm:"column:total_student"`
	TypeClass          string `json:"type_class" gorm:"column:type_class"`
}

func (ClassCurriculumnDetail) TableName() string {
	return "class_curriculum_details"
}
