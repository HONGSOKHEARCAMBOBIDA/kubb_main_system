package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type AcademicStreamService interface {
	GetAcademicStream(ctx context.Context) ([]model.AcademicStream, error)
}

type academicstreamservice struct {
	db *gorm.DB
}

func NewAcademicStreamService() AcademicStreamService {
	return &academicstreamservice{
		db: config.DB,
	}
}

func (s *academicstreamservice) GetAcademicStream(ctx context.Context) ([]model.AcademicStream, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.AcademicStream
	err := s.db.WithContext(ctx).Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch filingcabinet", err)
	}
	return data, nil
}
