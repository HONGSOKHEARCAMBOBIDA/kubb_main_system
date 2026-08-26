package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/utils"

	"gorm.io/gorm"
)

type ScheduleService interface {
	CreateScheduleService(ctx context.Context, input request.ScheduleRequest) error
}

type scheduleservice struct {
	db *gorm.DB
}

func NewScheduleService() ScheduleService {
	return &scheduleservice{
		db: config.DB,
	}
}

func (s *scheduleservice) CreateScheduleService(ctx context.Context, input request.ScheduleRequest) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Schedule{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		TeacherRateID:  input.TeacherRateID,
		ScheduleDate:   input.ScheduleDate,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		TotalTeachHour: input.TotalTeachHour,
		Description:    input.Description,
		Active:         true,
		RoomID:         input.RoomID,
		Status:         model.ScheduleStatusActive,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create teacher rate", nil)
		}
		return nil
	})
	return err
}
