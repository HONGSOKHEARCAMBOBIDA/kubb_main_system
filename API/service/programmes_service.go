package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type ProgramesService interface {
	GetProgrammes(ctx context.Context) ([]model.Programmes, error)
}

type programmesservice struct {
	db *gorm.DB
}

func NewProgrammesService() ProgramesService {
	return &programmesservice{
		db: config.DB,
	}
}

func (s *programmesservice) GetProgrammes(ctx context.Context) ([]model.Programmes, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Programmes
	err := s.db.WithContext(ctx).Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch filingcabinet", err)
	}
	return data, nil
}
