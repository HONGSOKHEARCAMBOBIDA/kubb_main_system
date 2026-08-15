package model

import "mysql/model/base"

type Province struct {
	base.ModelBase
	NameKh string `json:"name_kh" gorm:"column:name_kh"`
	NameEn string `json:"name_en" gorm:"column:name_en"`
}

type District struct {
	base.ModelBase
	NameKh     string `json:"name_kh" gorm:"column:name_kh"`
	NameEn     string `json:"name_en" gorm:"column:name_en"`
	ProvinceID int    `json:"province_id" gorm:"column:province_id"`
}

type Communce struct {
	base.ModelBase
	NameKh     string `json:"name_kh" gorm:"column:name_kh"`
	NameEn     string `json:"name_en" gorm:"column:name_en"`
	DistrictID int    `json:"district_id" gorm:"column:district_id"`
}

type Village struct {
	base.ModelBase
	NameKh     string `json:"name_kh" gorm:"column:name_kh"`
	NameEn     string `json:"name_en" gorm:"column:name_en"`
	CommunceID int    `json:"commune_id" gorm:"column:commune_id"`
}
