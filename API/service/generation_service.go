package service

import (
	"context"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GenerationService interface {
	GetGeneration(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.GenerationResponse, *model.PaginationMetadata, error)
	CreateGeneration(ctx context.Context, input request.GenerationRequestCreate) error
	UpdateGeneration(ctx context.Context, id string, input request.GenerationRequestUpdate) error
	Toggle(ctx context.Context, id string) error
	GetGenerationByAcademic(ctx context.Context, academicID int) ([]response.GenerationResponseByAcademic, error)
}

type generationservice struct {
	db *gorm.DB
}

func NewGenerationService() GenerationService {
	return &generationservice{
		db: config.DB,
	}
}

func (s *generationservice) GetGeneration(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.GenerationResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.GenerationResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("generations g").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["academic_id"]; ok && v != "" {
			tx = tx.Where("a.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count terms: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).Select(`
		g.id AS id,
		g.uuid AS uuid,
		a.id AS academic_id,
		a.name AS academic_name,
		g.code AS code,
		g.name AS name,
		g.start_date AS start_date,
		g.end_date AS end_date,
		g.description AS description,
		g.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *generationservice) GetGenerationByAcademic(ctx context.Context, academicID int) ([]response.GenerationResponseByAcademic, error) {
	if academicID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.GenerationResponseByAcademic

	err := s.db.WithContext(ctx).
		Table("generations g").
		Select(`
			g.id AS id,
			g.name AS name
		`).
		Where("g.academic_id = ?", academicID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch generations by academic: %w", err)
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
	return utils.ToggleStatus[model.Generation](ctx, s.db, id)
}
