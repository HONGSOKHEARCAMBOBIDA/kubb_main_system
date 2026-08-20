package response

import (
	"mysql/model"
	"mysql/model/base"
)

type InstallmentResponse struct {
	base.ModelBase
	base.UUIDBase
	FeeID           int                     `json:"fee_id" gorm:"column:fee_id"`
	SequenceNO      int                     `json:"sequence_no" gorm:"column:sequence_no"`
	DueDate         string                  `json:"due_date" gorm:"column:due_date"`
	Amount          float64                 `json:"amount" gorm:"column:amount"`
	Status          model.InstallmentStatus `json:"status"`
	InvoiceID       *int                    `json:"invoice_id" gorm:"column:invoice_id"`
	InvoiceResposne InvoiceResposne         `json:"invoice" gorm:"-"`
}
