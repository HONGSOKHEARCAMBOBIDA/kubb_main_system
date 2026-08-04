package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/utils"
	"strings"

	"gorm.io/gorm"
)

type AcademicService interface {
	GetAcademic(ctx context.Context) ([]model.Academic, error)
	CreateAcademic(ctx context.Context, input request.AcademicRequestCreate) error
	UpdateAcademic(ctx context.Context, id string, input request.AcademicRequestUpdate) error
	Toggle(ctx context.Context, id string) error
}

type academicservice struct {
	db *gorm.DB
}

func NewAcademicService() AcademicService {
	return &academicservice{
		db: config.DB,
	}
}

func (s *academicservice) GetAcademic(ctx context.Context) ([]model.Academic, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Academic
	err := s.db.WithContext(ctx).Order("id DESC").Find(&data).Error
	for i := range data {
		data[i].StartDate = helper.FormatDate(data[i].StartDate)
		data[i].EndDate = helper.FormatDatePtr(data[i].EndDate)
	}
	if err != nil {
		return nil, apperror.Internal("failed to fetch filingcabinet", err)
	}
	return data, nil
}

func (s *academicservice) CreateAcademic(ctx context.Context, input request.AcademicRequestCreate) error {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "academic code is required", nil)
	}
	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "academic name is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Academic{
		UUID:        helper.GenerateUUID(),
		Code:        code,
		Name:        name,
		StartDate:   input.StartDate,
		EndDate:     nil,
		Description: input.Description,
		Active:      true,
	}

	err := s.db.WithContext(ctx).Create(&newdata).Error
	if err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}

func (s *academicservice) UpdateAcademic(ctx context.Context, id string, input request.AcademicRequestUpdate) error {
	// if id != "" {
	// 	return apperror.Invalid("id must be a positive integer", nil)
	// }

	updates := map[string]interface{}{}

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return apperror.Invalid("code cannot be empty", nil)
		}
		updates["code"] = code
	}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.Invalid("name cannot be empty", nil)
		}
		updates["name"] = name
	}
	if input.StartDate != nil {
		updates["start_date"] = *input.StartDate
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}

	if len(updates) == 0 {
		return apperror.Invalid("no fields provided to update", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).
		Model(&model.Academic{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "update")
	}

	return nil
}

func (s *academicservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Academic](ctx, s.db, id)
}
