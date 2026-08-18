package request

type FeeRequestCreate struct {
	Date     string  `json:"date" gorm:"column:date"`
	Amount   float64 `json:"amount" gorm:"column:amount"`
	Discount float64 `json:"discount" gorm:"column:discount"`
	Total    float64 `json:"total" gorm:"column:total"`
	Active   bool    `json:"active" gorm:"column:discount"`
}
