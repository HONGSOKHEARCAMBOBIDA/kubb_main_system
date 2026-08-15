package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type LocationService interface {
	GetProvince(ctx context.Context) ([]model.Province, error)
	GetDistrict(ctx context.Context, id int) ([]model.District, error)
	GetCommunce(ctx context.Context, id int) ([]model.Communce, error)
	GetVillage(ctx context.Context, id int) ([]model.Village, error)
}

type locationservice struct {
	db *gorm.DB
}

func NewLocationService() LocationService {
	return &locationservice{
		db: config.DB,
	}
}

func (s *locationservice) GetProvince(ctx context.Context) ([]model.Province, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Province
	err := s.db.WithContext(ctx).Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch ", err)
	}
	return data, nil
}

func (s *locationservice) GetDistrict(ctx context.Context, id int) ([]model.District, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.District
	err := s.db.WithContext(ctx).Where("province_id =?", id).Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch ", err)
	}
	return data, nil
}

func (s *locationservice) GetCommunce(ctx context.Context, id int) ([]model.Communce, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Communce
	err := s.db.WithContext(ctx).Where("district_id =?", id).Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch ", err)
	}
	return data, nil
}

func (s *locationservice) GetVillage(ctx context.Context, id int) ([]model.Village, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Village
	err := s.db.WithContext(ctx).Where("commune_id =?", id).Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch ", err)
	}
	return data, nil
}
