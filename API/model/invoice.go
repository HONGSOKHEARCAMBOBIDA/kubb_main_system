package model

import "mysql/model/base"

type Invoice struct {
	base.ModelBase
	base.UUIDBase
	FeeID            int     `json:"fee_id" gorm:"column:fee_id"`
	Code             string  `json:"code" gorm:"column:code"`
	InvoiceDate      string  `json:"invoice_date" gorm:"column:invoice_date"`
	DueDate          string  `json:"due_date" gorm:"column:due_date"`
	Total            float64 `json:"total" gorm:"column:total"`
	Discount         float64 `json:"discount" gorm:"column:discount"`
	Tax              float64 `json:"tax" gorm:"column:tax"`
	GrantTotal       float64 `json:"grant_total" gorm:"column:grant_total"`
	MessageOnInvoice string  `json:"message_on_invoice" gorm:"column:message_on_invoice"`
	Description      string  `json:"description" gorm:"column:description"`
	Active           bool    `json:"active" gorm:"column:active"`
}

func (Invoice) TableName() string {
	return "invoices"
}
