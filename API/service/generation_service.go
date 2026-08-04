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
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GenerationService interface {
	GetGeneration(ctx context.Context) ([]model.Generation, error)
	CreateGeneration(ctx context.Context, input request.GenerationRequestCreate) error
	UpdateGeneration(ctx context.Context, id string, input request.GenerationRequestUpdate) error
	Toggle(ctx context.Context, id string) error
}

type generationservice struct {
	db *gorm.DB
}

func NewGenerationService() GenerationService {
	return &generationservice{
		db: config.DB,
	}
}

func (s *generationservice) GetGeneration(ctx context.Context) ([]model.Generation, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.Generation
	err := s.db.WithContext(ctx).Preload("Academic").Find(&data).Error
	for i := range data {
		data[i].StartDate = helper.FormatDate(data[i].StartDate)
		data[i].EndDate = helper.FormatDatePtr(data[i].EndDate)
	}
	if err != nil {
		return nil, apperror.Internal("failed to fetch filingcabinet", err)
	}
	return data, nil
}

func (s *generationservice) CreateGeneration(ctx context.Context, input request.GenerationRequestCreate) error {

	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation code is required", nil)
	}

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation name is required", nil)
	}

	if input.AcademicID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	nextIndex := 1

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lastGeneration model.Generation
		err := tx.
			Where("academic_id = ?", input.AcademicID).
			Order("id DESC").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastGeneration).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			nextIndex = 1
		case err != nil:
			return helper.MapAcademicError(err, "LOOKUP_LAST_GENERATION")
		default:
			nextIndex = lastGeneration.Index + 1
		}

		newdata := model.Generation{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			AcademicID:  input.AcademicID,
			Code:        code,
			Name:        name,
			Index:       nextIndex,
			StartDate:   input.StartDate,
			EndDate:     nil,
			Description: input.Description,
			Active:      true,
		}

		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapAcademicError(err, "CREATE")
		}
		return nil
	})

	return err
}

func (s *generationservice) UpdateGeneration(ctx context.Context, id string, input request.GenerationRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.AcademicID != nil {
		updates["academic_id"] = *input.AcademicID
	}

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return apperror.New(apperror.CodeInvalidInput, "generation code is required", nil)
		}
		updates["code"] = code
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "generation name is required", nil)
		}
		updates["name"] = name
	}

	if input.StartDate != nil {
		updates["start_date"] = *input.StartDate
	}

	if input.EndDate != nil {
		updates["end_date"] = *input.EndDate
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
		Model(&model.Generation{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "update")
	}

	if result.RowsAffected == 0 {
		return apperror.NotFound("generation not found", nil)
	}

	return nil
}

func (s *generationservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Academic](ctx, s.db, id)
}
