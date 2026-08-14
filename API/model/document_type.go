package model

import "mysql/model/base"

type DocumentType struct {
	base.ModelBase
	Code   string `json:"code" gorm:"column:code"`
	NameKh string `json:"name_kh" gorm:"column:name_kh"`
	NameEn string `json:"name_en" gorm:"column:name_en"`
}

func (DocumentType) TableName() string {
	return "document_types"
}
