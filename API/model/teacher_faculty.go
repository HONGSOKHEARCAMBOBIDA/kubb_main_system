package model

import "mysql/model/base"

type TeacherFaculty struct {
	base.ModelBase
	TeacherID int `json:"teacher_id" gorm:"column:teacher_id"`
	FacultyID int `json:"faculty_id" gorm:"column:faculty_id"`
}

func (TeacherFaculty) TableName() string {
	return "teacher_faculty"
}
