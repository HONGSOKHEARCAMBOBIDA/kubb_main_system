package model

import "mysql/model/base"

type ClassOffering struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnDetailID  int                 `json:"class_curriculum_detail_id" gorm:"column:class_curriculum_detail_id"`
	SubjectID                 int                 `json:"subject_id" gorm:"column:subject_id"`
	Credit                    int                 `json:"credit" gorm:"column:credit"`
	PassingScore              float64             `json:"passing_score" gorm:"column:passing_score"`
	TotalHour                 float64             `json:"total_hour" gorm:"column:total_hour"`
	Status                    StatusClassOffering `gorm:"type:enum('open','closed');not null"`
	TotalAttendanceForRexam   float64             `json:"total_attendance_for_rexam" gorm:"column:total_attendance_for_rexam"`
	TotalAttendanceForRelearn float64             `json:"total_attendance_for_relearn" gorm:"column:total_attendance_for_relearn"`
	Description               string              `json:"description" gorm:"column:description"`
}
