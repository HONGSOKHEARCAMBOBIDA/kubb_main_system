package model

type InstallmentStatus string

const (
	InstallmentStatusPending  InstallmentStatus = "pending"
	InstallmentStatusInvoiced InstallmentStatus = "invoiced"
	InstallmentStatusPaid     InstallmentStatus = "paid"
	InstallmentStatusOverdue  InstallmentStatus = "overdue"
)
