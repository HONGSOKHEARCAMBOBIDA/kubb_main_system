package response

import "mysql/model/base"

type FeeResponse struct {
	base.ModelBase
	base.UUIDBase
	EnrollmentID        int                   `json:"enrollment_id" gorm:"column:enrollment_id"`
	Date                string                `json:"date" gorm:"column:date"`
	Amount              float64               `json:"amount" gorm:"column:amount"`
	Discount            float64               `json:"discount" gorm:"column:discount"`
	Total               float64               `json:"total" gorm:"column:total"`
	Active              bool                  `json:"active" gorm:"column:active"`
	InvoiceResposne     []InvoiceResposne     `json:"invoice" gorm:"-"`
	InstallmentResponse []InstallmentResponse `json:"installment" gorm:"-"`
}
