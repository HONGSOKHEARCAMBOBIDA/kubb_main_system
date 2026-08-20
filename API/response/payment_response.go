package response

import "mysql/model/base"

type PaymentResposen struct {
	base.ModelBase
	base.UUIDBase
	InvoiceID   int     `json:"invoice_id" gorm:"column:invoice_id"`
	Code        string  `json:"code" gorm:"column:code"`
	Date        string  `json:"date" gorm:"column:date"`
	Amount      float64 `json:"amount" gorm:"column:amount"`
	Reference   string  `json:"reference" gorm:"column:reference"`
	Method      string  `json:"method" gorm:"column:method"`
	Description string  `json:"description" gorm:"column:description"`
	Active      bool    `json:"active" gorm:"column:active"`
}
