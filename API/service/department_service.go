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

type DepartmentService interface {
	GetDepartment(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.DepartmentResponse, *model.PaginationMetadata, error)
	CreateDepartment(ctx context.Context, input request.DepartmentRequestCreate) error
	UpdateDepartment(ctx context.Context, id string, input request.DepartmentRequestUpdate) error
	Toggle(ctx context.Context, id string) error
	GetDepartmentByFaculty(ctx context.Context, facultyID int) ([]response.DepartmentResponseByFaculty, error)
}

type departmentservice struct {
	db *gorm.DB
}

func NewDepartmentService() DepartmentService {
	return &departmentservice{
		db: config.DB,
	}
}

func (s *departmentservice) GetDepartment(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.DepartmentResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.DepartmentResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("departments d").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["faculty_id"]; ok && v != "" {
			tx = tx.Where("f.id = ?", v)
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
		d.id AS id,
		d.uuid AS uuid,
		f.id AS faculty_id,
		f.name AS faculty_name,
		f.code AS faculty_code,
		p.id AS programme_id,
		p.name AS programme_name,
		d.code AS code,
		d.name AS name,
		d.description AS description,
		d.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *departmentservice) CreateDepartment(ctx context.Context, input request.DepartmentRequestCreate) error {
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

	if input.FacultyID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "programme_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Department{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		FacultyID:   input.FacultyID,
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

func (s *departmentservice) UpdateDepartment(ctx context.Context, id string, input request.DepartmentRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.FacultyID != nil {
		updates["faculty_id"] = *input.FacultyID
	}

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return apperror.New(apperror.CodeInvalidInput, "department code is required", nil)
		}
		updates["code"] = code
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "department name is required", nil)
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
		Model(&model.Department{}).
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

func (s *departmentservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Department](ctx, s.db, id)
}

func (s *departmentservice) GetDepartmentByFaculty(ctx context.Context, facultyID int) ([]response.DepartmentResponseByFaculty, error) {
	if facultyID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "programme  is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.DepartmentResponseByFaculty

	err := s.db.WithContext(ctx).
		Table("departments d").
		Select(`
			d.id AS id,
			d.name AS name
		`).
		Where("d.faculty_id = ?", facultyID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch faculty by academic: %w", err)
	}

	return data, nil
}
