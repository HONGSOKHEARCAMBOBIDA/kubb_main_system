package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/utils"

	"gorm.io/gorm"
)

type SchoolRoomService interface {
	GetSchoolRoom(ctx context.Context) ([]model.SchoolRoom, error)
}

type schoolroomservice struct {
	db *gorm.DB
}

func NewSchoolRoomService() SchoolRoomService {
	return &schoolroomservice{
		db: config.DB,
	}
}

func (s *schoolroomservice) GetSchoolRoom(ctx context.Context) ([]model.SchoolRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.SchoolRoom
	err := s.db.WithContext(ctx).Preload("Floor.Building.Campuse.School").Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("faild to fetch data", err)
	}
	return data, nil
}
