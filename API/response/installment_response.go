package response

import (
	"mysql/model"
	"mysql/model/base"
)

type InstallmentResponse struct {
	base.ModelBase
	base.UUIDBase
	FeeID                   int                     `json:"fee_id" gorm:"column:fee_id"`
	SequenceNO              int                     `json:"sequence_no" gorm:"column:sequence_no"`
	DueDate                 string                  `json:"due_date" gorm:"column:due_date"`
	Amount                  float64                 `json:"amount" gorm:"column:amount"`
	Status                  model.InstallmentStatus `json:"status"`
	InvoiceID               *int                    `json:"invoice_id" gorm:"column:invoice_id"`
	InvoiceCode             string                  `json:"invoice_code"`
	InvoiceDate             string                  `json:"invoice_date" gorm:"column:invoice_date"`
	InvoiceDueDate          string                  `json:"invoice_due_date"`
	InvoiceTotal            float64                 `json:"invoice_total"`
	InvoiceTax              float64                 `json:"invoice_tax"`
	InvoiceGranTotal        float64                 `json:"invoice_grant_total" gorm:"column:invoice_grant_total"`
	InvoiceMessageOnInvoice string                  `json:"message_on_invoice" gorm:"column:message_on_invoice"`
	PaymentID               int                     `json:"payment_id"`
	PaymentCode             string                  `json:"payment_code"`
	PaymentReference        string                  `json:"payment_reference"`
	PaymentMethod           string                  `json:"payment_method"`
}
