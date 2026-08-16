package response

import "mysql/model/base"

type StudentEducationResponse struct {
	base.ModelBase
	base.UUIDBase
	StudentID       int    `json:"student_id" gorm:"column:student_id"`
	Level           string `json:"level" gorm:"column:level"`
	SchoolName      string `json:"school_name" gorm:"column:school_name"`
	VillageID       int    `json:"village_id" gorm:"column:village_id"`
	VillageNameKh   string `json:"villlage_name_kh" gorm:"column:villlage_name_kh"`
	CommunceID      int    `json:"communce_id"`
	CommunceName    string `json:"communce_name"`
	DistrictID      int    `json:"district_id"`
	DistrictName    string `json:"distirct_name"`
	ProvinceID      int    `json:"province_id"`
	ProvinceName    string `json:"province_name"`
	StartDate       string `json:"start_date" gorm:"column:start_date"`
	EndDate         string `json:"end_date" gorm:"column:end_date"`
	CertificateDate string `json:"cerificate_date" gorm:"column:cerificate_date"`
	Score           string `json:"score" gorm:"column:score"`
	Gpa             string `json:"gpa" gorm:"column:gpa"`
	Grade           string `json:"grade" gorm:"column:grade"`
}
