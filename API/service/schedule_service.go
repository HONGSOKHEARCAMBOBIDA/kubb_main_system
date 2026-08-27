package service

import (
	"context"
	"errors"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleService interface {
	CreateScheduleService(ctx context.Context, input request.ScheduleRequest) error
	UpdateScheduleService(ctx context.Context, uuid string, input request.ScheduleRequestUpdate, userID int) error
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
		VerifyBy:       nil,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create teacher rate", nil)
		}
		return nil
	})
	return err
}

func (s *scheduleservice) UpdateScheduleService(ctx context.Context, uuid string, input request.ScheduleRequestUpdate, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	if input.TotalTeachHour < 0 {
		return apperror.New(apperror.CodeInvalidInput, "total teach hour must be non-negative", nil)
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var schedule model.Schedule

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uuid = ?", uuid).
			First(&schedule).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "schedule not found", err)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch schedule", err)
		}

		updates := map[string]interface{}{
			"total_teach_hours": input.TotalTeachHour,
			"status":            input.Status,
			"verify_by":         userID,
		}

		if err := tx.Model(&schedule).Updates(updates).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update schedule", err)
		}

		return nil
	})

	return err
}
