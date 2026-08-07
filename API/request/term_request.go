package request

type TermRequestCreate struct {
	GenerationID int    `json:"generation_id" gorm:"column:generation_id"`
	Code         string `json:"code" gorm:"column:code"`
	Name         string `json:"name" gorm:"column:name"`
	Index        int    `json:"index" gorm:"column:index"`
	StartDate    string `json:"start_date" gorm:"column:start_date"`
	EndDate      string `json:"end_date" gorm:"column:end_date"`
	Description  string `json:"description" gorm:"column:description"`
}

type TermRequestUpdate struct {
	GenerationID *int    `json:"generation_id" gorm:"column:generation_id"`
	Code         *string `json:"code" gorm:"column:code"`
	Name         *string `json:"name" gorm:"column:name"`
	StartDate    *string `json:"start_date" gorm:"column:start_date"`
	EndDate      *string `json:"end_date" gorm:"column:end_date"`
	Description  *string `json:"description" gorm:"column:description"`
}
