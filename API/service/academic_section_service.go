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

type AcademicSectionService interface {
	CreateAcademicSection(ctx context.Context, input request.AcademicSectionRequestCreate) error
	GetAcademicSection(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AcademicSectionResponse, *model.PaginationMetadata, error)
	Toggle(ctx context.Context, id string) error
	UpdateAcademicSection(ctx context.Context, id string, input request.AcademicSectionRequestUpdate) error
}

type academicsectionservice struct {
	db *gorm.DB
}

func NewAcademicSectionService() AcademicSectionService {
	return &academicsectionservice{
		db: config.DB,
	}
}

func (s *academicsectionservice) CreateAcademicSection(ctx context.Context, input request.AcademicSectionRequestCreate) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "subject name is required", nil)
	}
	if input.MajorID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "major_id is required", nil)
	}
	if input.ShiftID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "major_id is required", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	newdata := model.AcademicSection{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		MajorID:     input.MajorID,
		ShiftID:     input.ShiftID,
		Name:        name,
		Description: input.Description,
		Type:        input.Type,
		Active:      true,
	}
	if err := s.db.WithContext(ctx).Create(&newdata).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}

func (s *academicsectionservice) GetAcademicSection(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AcademicSectionResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.AcademicSectionResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("academic_sections ash").
			Joins("LEFT JOIN majors m ON m.id = ash.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
			Joins("LEFT JOIN academic_shifts asf ON asf.id = ash.shift_id").
			Joins("LEFT JOIN academics a ON a.id = asf.academic_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["programme_id"]; ok && v != "" {
			tx = tx.Where("p.id = ?", v)
		}
		if v, ok := filter["faculty_id"]; ok && v != "" {
			tx = tx.Where("f.id = ?", v)
		}
		if v, ok := filter["department_id"]; ok && v != "" {
			tx = tx.Where("d.id = ?", v)
		}
		if v, ok := filter["major_id"]; ok && v != "" {
			tx = tx.Where("m.id = ?", v)
		}
		if v, ok := filter["shift_id"]; ok && v != "" {
			tx = tx.Where("asf.id = ?", v)
		}
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
		ash.type AS type,
		ash.active AS active,
		m.id AS major_id,
		m.code AS major_code,
		m.name AS major_name,
		p.id AS programme_id,
		p.name AS programme_name,
		d.id AS department_id,
		f.id AS faculty_id,
		asf.id AS shift_id,
		asf.name AS shift_name,
		a.id AS academic_id,
		a.code AS academic_code
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *academicsectionservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.AcademicSection](ctx, s.db, id)
}

func (s *academicsectionservice) UpdateAcademicSection(ctx context.Context, id string, input request.AcademicSectionRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "subject id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	updates := map[string]interface{}{}

	if input.MajorID != nil {
		if *input.MajorID == 0 {
			return apperror.New(apperror.CodeInvalidInput, "major_id cannot be zero", nil)
		}
		updates["major_id"] = *input.MajorID
	}

	if input.ShiftID != nil {
		if *input.ShiftID == 0 {
			return apperror.New(apperror.CodeInvalidInput, "major_id cannot be zero", nil)
		}
		updates["shift_id"] = *input.ShiftID
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "subject name cannot be empty", nil)
		}
		updates["name"] = name
	}

	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}

	if input.Type != nil {
		updates["type"] = *input.Type
	}

	if len(updates) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "no fields provided to update", nil)
	}

	result := s.db.WithContext(ctx).
		Model(&model.AcademicSection{}).
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
