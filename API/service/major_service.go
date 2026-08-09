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
	"strings"

	"gorm.io/gorm"
)

type MajorService interface {
	CreateMajor(ctx context.Context, input request.MajorRequestCreate) error
}

type majorservice struct {
	db *gorm.DB
}

func NewMajorService() MajorService {
	return &majorservice{
		db: config.DB,
	}
}

func (s *majorservice) CreateMajor(ctx context.Context, input request.MajorRequestCreate) error {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "major code is required", nil)
	}

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "major name is required", nil)
	}

	if input.DepartmentID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "department_id is required", nil)
	}

	if input.DurationPeriod <= 0 {
		return apperror.New(apperror.CodeInvalidInput, "duration_period must be greater than 0", nil)
	}

	durationInterval := input.DurationInterval
	if durationInterval == "" {
		durationInterval = model.DurationIntervalYear
	}
	switch durationInterval {
	case model.DurationIntervalYear,
		model.DurationIntervalMonth,
		model.DurationIntervalWeek,
		model.DurationIntervalDay:
		// ok
	default:
		return apperror.New(apperror.CodeInvalidInput, "invalid duration_interval, must be one of year, month, week, day", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Major{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		DepartmentID:     int(input.DepartmentID),
		Code:             code,
		Name:             name,
		DurationPeriod:   input.DurationPeriod,
		DurationInterval: durationInterval,
		Description:      description,
		Active:           true,
	}

	if err := s.db.WithContext(ctx).Create(&newdata).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}
