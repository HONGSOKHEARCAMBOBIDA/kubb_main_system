package service

import (
	"context"
	"fmt"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/utils"
	"time"

	"gorm.io/gorm"
)

type InvoiceService interface {
	CreateInvoice(ctx context.Context, input request.InvoiceRequestCreate) error
}

type invoiceservice struct {
	db *gorm.DB
}

func NewInvoiceService() InvoiceService {
	return &invoiceservice{
		db: config.DB,
	}
}

func (s *invoiceservice) CreateInvoice(ctx context.Context, input request.InvoiceRequestCreate) error {
	if input.FeeID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "fee id is required", nil)
	}
	if input.Total <= 0 {
		return apperror.New(apperror.CodeInvalidInput, "total is required", nil)
	}
	if input.GrantTotal <= 0 {
		return apperror.New(apperror.CodeInvalidInput, "grant total is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newinvoice := model.Invoice{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		FeeID:            input.FeeID,
		InvoiceDate:      input.InvoiceDate,
		DueDate:          input.DueDate,
		Total:            input.Total,
		Discount:         input.Discount,
		Tax:              input.Tax,
		GrantTotal:       input.GrantTotal,
		MessageOnInvoice: input.MessageOnInvoice,
		Active:           true,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Create(&newinvoice)
		newinvoice.Code = helper.GenerateCode("INVOICE", uint(newinvoice.ID))
		tx.Save(&newinvoice)

		newpayment := model.Payment{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			InvoiceID: newinvoice.ID,
			Date:      time.Now().Format("2006-01-02"),
			Amount:    input.GrantTotal,
			Reference: input.Reference,
			Method:    input.Method,
			Active:    true,
		}
		tx.Create(&newpayment)

		newpayment.Code = helper.GenerateCode("PAYMENT", uint(newpayment.ID))
		tx.Save(&newpayment)

		updates := map[string]interface{}{}

		updates["invoice_id"] = newinvoice.ID
		updates["status"] = model.InstallmentStatusPaid
		result := tx.Model(&model.Installment{}).Where("uuid = ?", input.InstallmentUUID).Updates(updates)
		if result.Error != nil {
			return helper.MapAcademicError(result.Error, "update")
		}
		fmt.Println("rows affected:", result.RowsAffected, "installment uuid:", input.InstallmentUUID)
		if result.RowsAffected == 0 {
			return apperror.New(apperror.CodeNotFound, "installment not found", nil)
		}
		return nil
	})

	return err
}
