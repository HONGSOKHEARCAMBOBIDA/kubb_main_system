package model

import (
	"mysql/model/base"
)

type StudentFamily struct {
	base.ModelBase
	base.UUIDBase

	StudentID int `json:"student_id" gorm:"column:student_id"`

	FatherName        string `json:"father_name" gorm:"column:father_name"`
	FatherEnglishName string `json:"father_english_name" gorm:"column:father_english_name"`
	FatherAge         int    `json:"father_age" gorm:"column:father_age"`
	FatherIsAlive     bool   `json:"father_is_alive" gorm:"column:father_is_alive"`
	FatherPhoneNumber string `json:"father_phone_number" gorm:"column:father_phone_number"`
	FatherOccupation  string `json:"father_occupation" gorm:"column:father_occupation"`
	FatherWorkplace   string `json:"father_workplace" gorm:"column:father_workplace"`

	MotherName        string `json:"mother_name" gorm:"column:mother_name"`
	MotherEnglishName string `json:"mother_english_name" gorm:"column:mother_english_name"`
	MotherAge         int    `json:"mother_age" gorm:"column:mother_age"`
	MotherIsAlive     bool   `json:"mother_is_alive" gorm:"column:mother_is_alive"`
	MotherPhoneNumber string `json:"mother_phone_number" gorm:"column:mother_phone_number"`
	MotherOccupation  string `json:"mother_occupation" gorm:"column:mother_occupation"`
	MotherWorkplace   string `json:"mother_workplace" gorm:"column:mother_workplace"`
}
