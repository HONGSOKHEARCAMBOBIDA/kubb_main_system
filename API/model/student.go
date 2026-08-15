package model

import "mysql/model/base"

type Student struct {
	base.ModelBase
	base.UUIDBase
	GroupID          int    `json:"group_id" gorm:"column:group_id"`
	Code             string `json:"code" gorm:"column:code"`
	UserName         string `json:"username" gorm:"column:username"`
	Email            string `json:"email" gorm:"column:email"`
	Password         string `json:"password" gorm:"column:password"`
	NameKh           string `json:"name_kh" gorm:"column:name_kh"`
	NameEn           string `json:"name_en" gorm:"column:name_en"`
	DateOfBirth      string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender           string `json:"gender" gorm:"column:gender"`
	Nationality      string `json:"nationality" gorm:"column:nationality"`
	Phone            string `json:"phone" gorm:"column:phone"`
	Status           string `json:"status" gorm:"column:status"`
	VillageID        int    `json:"village_id" gorm:"column:village_id"`
	Occupation       string `json:"occupation" gorm:"column:occupation"`
	AcademicStreamID int    `json:"academic_stream_id" gorm:"column:academic_stream_id"`
	TelegramUsername string `json:"telegram_username" gorm:"column:telegram_username"`
}
