package service

import (
	"context"
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
)

type AcademicShiftService interface {
	GetAcademicShift(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AcademicShiftResponse, *model.PaginationMetadata, error)
	CreateAcademicShift(ctx context.Context, input request.AcademicShiftRequestCreate) error
	UpdateAcademicShift(ctx context.Context, id string, input request.AcademicShiftRequestUpdate) error
	Toggle(ctx context.Context, id string) error
	GetAcademicShiftByAcademic(ctx context.Context, academicID int) ([]response.AcademicShiftResponseByAcademic, error)
}

type academicshiftservice struct {
	db *gorm.DB
}

func NewAcademicShiftService() AcademicShiftService {
	return &academicshiftservice{
		db: config.DB,
	}
}

func (s *academicshiftservice) GetAcademicShift(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AcademicShiftResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.AcademicShiftResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("academic_shifts ash").
			Joins("LEFT JOIN academics a ON a.id = ash.academic_id")
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
		ash.id AS id,
		ash.uuid AS uuid,
		ash.name AS name,
		ash.description AS description,
		ash.active AS active,
		a.id AS academic_id,
		a.code AS academic_code,
		a.name AS academic_name
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *academicshiftservice) GetAcademicShiftByAcademic(ctx context.Context, academicID int) ([]response.AcademicShiftResponseByAcademic, error) {
	if academicID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.AcademicShiftResponseByAcademic

	err := s.db.WithContext(ctx).
		Table("academic_shifts ash").
		Select(`
			ash.id AS id,
			ash.name AS name
		`).
		Where("ash.academic_id = ?", academicID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch generations by academic: %w", err)
	}

	return data, nil
}

func (s *academicshiftservice) CreateAcademicShift(ctx context.Context, input request.AcademicShiftRequestCreate) error {
	name := strings.TrimSpace(input.Name)

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation code is required", nil)
	}

	if input.AcademicID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.AcademicShift{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		AcademicID:  input.AcademicID,
		Name:        name,
		Description: input.Description,
		Active:      true,
	}

	if err := s.db.WithContext(ctx).Create(&newdata).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}

func (s *academicshiftservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.AcademicShift](ctx, s.db, id)
}

func (s *academicshiftservice) UpdateAcademicShift(ctx context.Context, id string, input request.AcademicShiftRequestUpdate) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.AcademicID != nil {
		updates["academic_id"] = *input.AcademicID
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "generation name is required", nil)
		}
		updates["name"] = name
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
		Model(&model.AcademicShift{}).
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
