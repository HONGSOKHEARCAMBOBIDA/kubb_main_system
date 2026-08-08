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

type SemesterService interface {
	GetSemester(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.SemesternResponse, *model.PaginationMetadata, error)
	CreateSemester(ctx context.Context, input request.SemesterRequestCreate) error
	UpdateSemester(ctx context.Context, id string, input request.SemesterRequestUpdate) error
	Toggle(ctx context.Context, id string) error
	GetSemesterByAcademic(ctx context.Context, academicID int) ([]response.SemesterResponseByAcademic, error)
}

type semesterservice struct {
	db *gorm.DB
}

func NewSemesterService() SemesterService {
	return &semesterservice{
		db: config.DB,
	}
}

func (s *semesterservice) GetSemester(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.SemesternResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.SemesternResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("semesters s").
			Joins("LEFT JOIN academics a ON a.id = s.academic_id")
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
		s.id AS id,
		s.uuid AS uuid,
		a.id AS academic_id,
		a.name AS academic_name,
		s.code AS code,
		s.name AS name,
		s.start_date AS start_date,
		s.end_date AS end_date,
		s.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *semesterservice) GetSemesterByAcademic(ctx context.Context, academicID int) ([]response.SemesterResponseByAcademic, error) {
	if academicID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.SemesterResponseByAcademic

	err := s.db.WithContext(ctx).
		Table("semesters s").
		Select(`
			s.id AS id,
			s.name AS name
		`).
		Where("s.academic_id = ?", academicID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch generations by academic: %w", err)
	}

	return data, nil
}

func (s *semesterservice) CreateSemester(ctx context.Context, input request.SemesterRequestCreate) error {

	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "semester code is required", nil)
	}

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "semester name is required", nil)
	}

	if input.AcademicID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	nextIndex := 1

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lastSemester model.Semester
		err := tx.
			Where("academic_id = ?", input.AcademicID).
			Order("id DESC").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastSemester).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			nextIndex = 1
		case err != nil:
			return helper.MapAcademicError(err, "LOOKUP_LAST_SEMESTER")
		default:
			nextIndex = lastSemester.Index + 1
		}

		newdata := model.Semester{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			AcademicID: input.AcademicID,
			Code:       code,
			Name:       name,
			Index:      nextIndex,
			StartDate:  input.StartDate,
			EndDate:    input.EndDate,
			Active:     true,
		}

		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapAcademicError(err, "CREATE")
		}
		return nil
	})

	return err
}

func (s *semesterservice) UpdateSemester(ctx context.Context, id string, input request.SemesterRequestUpdate) error {
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

	if len(updates) == 0 {
		return apperror.Invalid("no fields provided to update", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).
		Model(&model.Semester{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "update")
	}

	if result.RowsAffected == 0 {
		return apperror.NotFound("semester not found", nil)
	}

	return nil
}

func (s *semesterservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Semester](ctx, s.db, id)
}
