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

type MajorService interface {
	CreateMajor(ctx context.Context, input request.MajorRequestCreate) error
	GetMajor(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.MajorResponse, *model.PaginationMetadata, error)
	GetMajorByDepartment(ctx context.Context, departmentID int) ([]response.MajorResponseByDepartment, error)
	Toggle(ctx context.Context, id string) error
	UpdateMajor(ctx context.Context, id string, input request.MajorRequestUpdate) error
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

func (s *majorservice) GetMajor(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.MajorResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.MajorResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("majors m").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id")
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
		m.id AS id,
		m.uuid AS uuid,
		m.code AS code ,
		m.name AS name,
		m.duration_period AS duration_period,
		m.duration_interval AS duration_interval,
		m.description AS description,
		m.active AS active,
		d.id AS department_id,
		d.name AS department_name,
		d.code AS department_code,
		f.id AS faculty_id,
		f.code AS faculty_code,
		f.name AS faculty_name,
		p.id AS programme_id,
		p.name AS programme_name
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *majorservice) GetMajorByDepartment(ctx context.Context, departmentID int) ([]response.MajorResponseByDepartment, error) {
	if departmentID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "department id  is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.MajorResponseByDepartment

	err := s.db.WithContext(ctx).
		Table("majors m").
		Select(`
			m.id AS id,
			m.name AS name
		`).
		Where("m.department_id = ?", departmentID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch major by department: %w", err)
	}

	return data, nil
}

func (s *majorservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Major](ctx, s.db, id)
}

func (s *majorservice) UpdateMajor(ctx context.Context, id string, input request.MajorRequestUpdate) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.DepartmentID != nil {
		if *input.DepartmentID == 0 {
			return apperror.New(apperror.CodeInvalidInput, "department_id is required", nil)
		}
		updates["department_id"] = int(*input.DepartmentID)
	}

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return apperror.New(apperror.CodeInvalidInput, "major code is required", nil)
		}
		updates["code"] = code
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "major name is required", nil)
		}
		updates["name"] = name
	}

	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}

	if input.DurationPeriod != nil {
		if *input.DurationPeriod <= 0 {
			return apperror.New(apperror.CodeInvalidInput, "duration_period must be greater than 0", nil)
		}
		updates["duration_period"] = *input.DurationPeriod
	}

	if input.DurationInterval != nil {
		durationInterval := *input.DurationInterval
		switch durationInterval {
		case model.DurationIntervalYear,
			model.DurationIntervalMonth,
			model.DurationIntervalWeek,
			model.DurationIntervalDay:
		default:
			return apperror.New(apperror.CodeInvalidInput, "invalid duration_interval, must be one of year, month, week, day", nil)
		}
		updates["duration_interval"] = durationInterval
	}

	if len(updates) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "no fields provided to update", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).
		Model(&model.Major{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "UPDATE")
	}
	if result.RowsAffected == 0 {
		return apperror.New(apperror.CodeNotFound, "major not found", nil)
	}

	return nil
}
