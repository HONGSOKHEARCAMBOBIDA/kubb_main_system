package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type SchoolOfficeService interface {
	GetSchoolOffice(ctx context.Context) ([]model.SchoolOffice, error)
}

type schoolofficeservice struct {
	db *gorm.DB
}

func NewSchoolOfficeService() SchoolOfficeService {
	return &schoolofficeservice{
		db: config.DB,
	}
}

func (s *schoolofficeservice) GetSchoolOffice(ctx context.Context) ([]model.SchoolOffice, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.SchoolOffice
	err := s.db.WithContext(ctx).Preload("Floor.Building.Campuse.School").Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("faild to fetch data", err)
	}
	return data, nil
}
