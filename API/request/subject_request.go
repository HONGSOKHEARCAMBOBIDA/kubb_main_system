package request

type SubjectRequestCreate struct {
	MajorID      int     `json:"major_id" gorm:"column:major_id"`
	Code         string  `json:"code" gorm:"column:code"`
	Name         string  `json:"name" gorm:"column:name"`
	Credit       float64 `json:"credit" gorm:"column:credit"`
	PassingScore float64 `json:"passing_score" gorm:"column:passing_score"`
	Description  string  `json:"description" gorm:"column:description"`
}

type SubjectRequestUpdate struct {
	MajorID      *int     `json:"major_id" gorm:"column:major_id"`
	Code         *string  `json:"code" gorm:"column:code"`
	Name         *string  `json:"name" gorm:"column:name"`
	Credit       *float64 `json:"credit" gorm:"column:credit"`
	PassingScore *float64 `json:"passing_score" gorm:"column:passing_score"`
	Description  *string  `json:"description" gorm:"column:description"`
}
