package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type BuildingService interface {
	GetBuilding(ctx context.Context) ([]model.Building, error)
}

type buildingsservice struct {
	db *gorm.DB
}

func NewBuildingService() BuildingService {
	return &buildingsservice{
		db: config.DB,
	}
}

func (s *buildingsservice) GetBuilding(ctx context.Context) ([]model.Building, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Building
	err := s.db.WithContext(ctx).Preload("Campuse").Preload("Campuse.School").Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("faild to fetch data", err)
	}
	return data, nil
}
