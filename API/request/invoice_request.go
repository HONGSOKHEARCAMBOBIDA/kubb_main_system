package request

type InvoiceRequestCreate struct {
	FeeID            int     `json:"fee_id" gorm:"column:fee_id"`
	InvoiceDate      string  `json:"invoice_date" gorm:"column:invoice_date"`
	DueDate          string  `json:"due_date" gorm:"column:due_date"`
	Total            float64 `json:"total" gorm:"column:total"`
	Discount         float64 `json:"discount" gorm:"column:discount"`
	Tax              float64 `json:"tax" gorm:"column:tax"`
	GrantTotal       float64 `json:"grant_total" gorm:"column:grant_total"`
	MessageOnInvoice string  `json:"message_on_invoice" gorm:"column:message_on_invoice"`
	Date             string  `json:"date" gorm:"column:date"`
	Amount           float64 `json:"amount" gorm:"column:amount"`
	Reference        string  `json:"reference" gorm:"column:reference"`
	Method           string  `json:"method" gorm:"column:method"`
}
