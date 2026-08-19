package model

import "mysql/model/base"

type Installment struct {
	base.ModelBase
	base.UUIDBase
	FeeID      int               `json:"fee_id" gorm:"column:fee_id"`
	SequenceNO int               `json:"sequence_no" gorm:"column:sequence_no"`
	DueDate    string            `json:"due_date" gorm:"column:due_date"`
	Amount     float64           `json:"amount" gorm:"column:amount"`
	Status     InstallmentStatus `gorm:"type:enum('pending','invoiced','paid','overdue');not null"`
	InvoiceID  *int              `json:"invoice_id" gorm:"column:invoice_id"`
}

func (Installment) TableName() string {
	return "installments"
}
