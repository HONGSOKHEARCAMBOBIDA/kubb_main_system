package service

import (
	"context"
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
		if err := tx.Create(&newinvoice).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create invoice", err)
		}

		newinvoice.Code = helper.GenerateCode("INVOICE", uint(newinvoice.ID))
		if err := tx.Save(&newinvoice).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update invoice code", err)
		}

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
		if err := tx.Create(&newpayment).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create payment", err)
		}

		newpayment.Code = helper.GenerateCode("PAYMENT", uint(newpayment.ID))
		if err := tx.Save(&newpayment).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update payment code", err)
		}

		return nil
	})

	return err
}
