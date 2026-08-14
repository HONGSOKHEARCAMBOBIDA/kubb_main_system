package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type DocumentTypeService interface {
	GetDocumentType(ctx context.Context) ([]model.DocumentType, error)
}

type documenttypeservice struct {
	db *gorm.DB
}

func NewDocumentTypeService() DocumentTypeService {
	return &documenttypeservice{
		db: config.DB,
	}
}

func (s *documenttypeservice) GetDocumentType(ctx context.Context) ([]model.DocumentType, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.DocumentType
	err := s.db.WithContext(ctx).Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch filingcabinet", err)
	}
	return data, nil
}
