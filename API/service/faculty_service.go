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

type FacultyService interface {
	GetFaculty(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.FacultyResponse, *model.PaginationMetadata, error)
	CreateFaculty(ctx context.Context, input request.FacultyRequestCreate) error
	UpdateFaculty(ctx context.Context, id string, input request.FacultyRequestUpdate) error
	Toggle(ctx context.Context, id string) error
	GetFacultyByProgrammes(ctx context.Context, programmeID int) ([]response.FacultyResponseByProgrammes, error)
}

type facultyservice struct {
	db *gorm.DB
}

func NewFacultyService() FacultyService {
	return &facultyservice{
		db: config.DB,
	}
}

func (s *facultyservice) GetFaculty(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.FacultyResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.FacultyResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("faculties f").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["programme_id"]; ok && v != "" {
			tx = tx.Where("p.id = ?", v)
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
		f.id AS id,
		f.uuid AS uuid,
		p.id AS programme_id,
		p.name AS programme_name,
		f.code AS code,
		f.name AS name,
		f.description AS description,
		f.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *facultyservice) CreateFaculty(ctx context.Context, input request.FacultyRequestCreate) error {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "faculty code is required", nil)
	}
	if len(code) > 20 {
		return apperror.New(apperror.CodeInvalidInput, "faculty code must not exceed 20 characters", nil)
	}

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "faculty name is required", nil)
	}
	if len(name) > 150 {
		return apperror.New(apperror.CodeInvalidInput, "faculty name must not exceed 150 characters", nil)
	}

	if input.ProgrammeID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "programme_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Faculty{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		ProgrammeID: input.ProgrammeID,
		Code:        code,
		Name:        name,
		Description: description,
		Active:      true,
	}

	if err := s.db.WithContext(ctx).Create(&newdata).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}

func (s *facultyservice) UpdateFaculty(ctx context.Context, id string, input request.FacultyRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.ProgrammeID != nil {
		updates["programme_id"] = *input.ProgrammeID
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

	if input.Description != nil {
		updates["description"] = *input.Description
	}

	if len(updates) == 0 {
		return apperror.Invalid("no fields provided to update", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).
		Model(&model.Faculty{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "update")
	}

	if result.RowsAffected == 0 {
		return apperror.NotFound("faculty not found", nil)
	}

	return nil
}

func (s *facultyservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Faculty](ctx, s.db, id)
}

func (s *facultyservice) GetFacultyByProgrammes(ctx context.Context, programmeID int) ([]response.FacultyResponseByProgrammes, error) {
	if programmeID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "programme  is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.FacultyResponseByProgrammes

	err := s.db.WithContext(ctx).
		Table("faculties f").
		Select(`
			f.id AS id,
			f.name AS name
		`).
		Where("f.programme_id = ?", programmeID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch faculty by academic: %w", err)
	}

	return data, nil
}
