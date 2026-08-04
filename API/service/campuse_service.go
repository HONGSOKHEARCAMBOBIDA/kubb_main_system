package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type CampuseService interface {
	GetCampuse(ctx context.Context) ([]model.Campuse, error)
}

type campuseservice struct {
	db *gorm.DB
}

func NewCampuseService() CampuseService {
	return &campuseservice{
		db: config.DB,
	}
}

func (s *campuseservice) GetCampuse(ctx context.Context) ([]model.Campuse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Campuse
	err := s.db.WithContext(ctx).Order("id DESC").Preload("School").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("faild to fetch data", err)
	}
	return data, nil
}
