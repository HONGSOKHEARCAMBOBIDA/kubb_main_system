package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type FloorService interface {
	GetFloor(ctx context.Context) ([]model.Floor, error)
}

type floorservice struct {
	db *gorm.DB
}

func NewFloorService() FloorService {
	return &floorservice{
		db: config.DB,
	}
}

func (s *floorservice) GetFloor(ctx context.Context) ([]model.Floor, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []model.Floor
	err := s.db.WithContext(ctx).
		Preload("Building.Campuse.School").
		Order("id DESC").
		Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch data", err)
	}
	return data, nil
}
