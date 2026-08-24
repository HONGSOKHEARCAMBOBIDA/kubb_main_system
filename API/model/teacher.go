package model

import "mysql/model/base"

type Teacher struct {
	base.ModelBase
	base.UUIDBase
	Code        string `json:"code" gorm:"column:code"`
	Email       string `json:"email" gorm:"column:email"`
	Password    string `json:"password" gorm:"column:password"`
	Name        string `json:"name" gorm:"column:name"`
	Dob         string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Pob         string `json:"place_of_birth" gorm:"column:place_of_birth"`
	Gender      string `json:"gender" gorm:"column:gender"`
	Nationality string `json:"nationality" gorm:"column:nationality"`
	Address     string `json:"address" gorm:"column:address"`
	Phone       string `json:"phone" gorm:"column:phone"`
	ProgrammeID int    `json:"programme_id" gorm:"column:programme_id"`
	FacultyID   int    `json:"faculty_id" gorm:"column:faculty_id"`
}
