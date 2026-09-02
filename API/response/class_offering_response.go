package response

import (
	"mysql/model"
	"mysql/model/base"
)

type ClassOfferingResponse struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnDetailID  int                       `json:"class_curriculum_detail_id" gorm:"column:class_curriculum_detail_id"`
	SubjectID                 int                       `json:"subject_id" gorm:"column:subject_id"`
	SubjectCode               string                    `json:"subject_code"`
	SubjectName               string                    `json:"subject_name"`
	SubjectGroupID            int                       `json:"subject_group_id"`
	SubjectGroupName          string                    `json:"subject_group_name"`
	Credit                    int                       `json:"credit" gorm:"column:credit"`
	PassingScore              float64                   `json:"passing_score" gorm:"column:passing_score"`
	TotalHour                 float64                   `json:"total_hour" gorm:"column:total_hour"`
	Status                    model.StatusClassOffering `json:"status"`
	TotalAttendanceForRexam   float64                   `json:"total_attendance_for_rexam" gorm:"column:total_attendance_for_rexam"`
	TotalAttendanceForRelearn float64                   `json:"total_attendance_for_relearn" gorm:"column:total_attendance_for_relearn"`
	Description               string                    `json:"description" gorm:"column:description"`
	StudentResponseSummary    []StudentResponseSummary  `json:"student" gorm:"-"`
}

type ClassOfferingResponseWithTeacherRate struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnDetailID  int                       `json:"class_curriculum_detail_id" gorm:"column:class_curriculum_detail_id"`
	SubjectID                 int                       `json:"subject_id" gorm:"column:subject_id"`
	SubjectCode               string                    `json:"subject_code"`
	SubjectName               string                    `json:"subject_name"`
	Credit                    int                       `json:"credit" gorm:"column:credit"`
	PassingScore              float64                   `json:"passing_score" gorm:"column:passing_score"`
	TotalHour                 float64                   `json:"total_hour" gorm:"column:total_hour"`
	RemainingHour             float64                   `json:"remaining_hour" gorm:"column:remaining_hour"`
	Status                    model.StatusClassOffering `json:"status"`
	TotalAttendanceForRexam   float64                   `json:"total_attendance_for_rexam" gorm:"column:total_attendance_for_rexam"`
	TotalAttendanceForRelearn float64                   `json:"total_attendance_for_relearn" gorm:"column:total_attendance_for_relearn"`
	Description               string                    `json:"description" gorm:"column:description"`
	TeacherRateID             int                       `json:"teacher_rate_id" gorm:"column:teacher_rate_id"`
	TeacherID                 int                       `json:"teacher_id"`
	TeacherName               string                    `json:"teacher_name"`
	TeacherGender             string                    `json:"teacher_gender"`
	HourlyRate                float64                   `json:"hourly_rate"`
	EffectiveDate             string                    `json:"effective_date"`
	Endate                    string                    `json:"end_date"`
	Active                    bool                      `json:"active"`
	ScheduleResponse          []ScheduleResponse        `json:"scheduel" gorm:"-"`
}
