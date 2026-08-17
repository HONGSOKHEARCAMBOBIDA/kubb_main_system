package response

import (
	"mysql/model"
	"mysql/model/base"
)

type StudentResponse struct {
	base.ModelBase
	base.UUIDBase
	GroupID            int                `json:"group_id" gorm:"column:group_id"`
	GroupCode          string             `json:"group_code"`
	GroupName          string             `json:"group_name"`
	DiscountType       model.DiscountType `json:"discount_type"`
	DiscountPercentage float64            `json:"discount_percentage" gorm:"column:discount_percentage"`
	DiscountAmount     float64            `json:"discount_amount" gorm:"column:discount_amount"`
	Code               string             `json:"code" gorm:"column:code"`
	UserName           string             `json:"username" gorm:"column:username"`
	Email              string             `json:"email" gorm:"column:email"`
	NameKh             string             `json:"name_kh" gorm:"column:name_kh"`
	NameEn             string             `json:"name_en" gorm:"column:name_en"`
	DateOfBirth        string             `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender             string             `json:"gender" gorm:"column:gender"`
	Nationality        string             `json:"nationality" gorm:"column:nationality"`
	Phone              string             `json:"phone" gorm:"column:phone"`
	Status             string             `json:"status" gorm:"column:status"`
	VillageID          int                `json:"village_id" gorm:"column:village_id"`
	VillageNameKh      string             `json:"villlage_name_kh" gorm:"column:villlage_name_kh"`
	CommunceID         int                `json:"communce_id" gorm:"column:communce_id"`
	CommunceName       string             `json:"communce_name"`
	DistrictID         int                `json:"district_id" gorm:"column:district_id"`
	DistrictName       string             `json:"distirct_name" gorm:"column:distirct_name"`
	ProvinceID         int                `json:"province_id"`
	ProvinceName       string             `json:"province_name"`
	Occupation         string             `json:"occupation" gorm:"column:occupation"`
	AcademicStreamID   int                `json:"academic_stream_id" gorm:"column:academic_stream_id"`
	AcademicStreamName string             `json:"academic_stream_name"`
	TelegramUsername   *string            `json:"telegram_username" gorm:"column:telegram_username"`

	StudentFamily    []model.StudentFamily      `json:"student_family" gorm:"-"`
	StudentEducation []StudentEducationResponse `json:"student_educations" gorm:"-"`
	StudentDocument  []StudentDocumentResponse  `json:"student_documents" gorm:"-"`
}
