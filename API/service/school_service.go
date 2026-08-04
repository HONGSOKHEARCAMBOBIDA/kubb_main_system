package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type SchoolService interface {
	GetSchool(ctx context.Context) ([]model.School, error)
}

type schoolservice struct {
	db *gorm.DB
}

func NewSchoolService() SchoolService {
	return &schoolservice{
		db: config.DB,
	}
}

func (s *schoolservice) GetSchool(ctx context.Context) ([]model.School, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.School
	err := s.db.WithContext(ctx).Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("faild to fetch data", err)
	}
	return data, nil
}
