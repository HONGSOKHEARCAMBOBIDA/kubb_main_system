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

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

type AcademicDegreeService interface {
	GetAcademicDegree(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AcademicDegreeResponse, *model.PaginationMetadata, error)
	GetAcademicDegreeByAcademicID(ctx context.Context, id int) ([]response.AcademicDegreeResponse, error)
	CreateAcademicDegree(ctx context.Context, input request.AcademicDegreeRequestCreate) error
	UpdateAcademicDegree(ctx context.Context, id string, input request.AcademicDegreeRequestUpdate) error
	Toggle(ctx context.Context, id string) error
}

type academicdegreeservice struct {
	db *gorm.DB
}

func NewAcademicDegreeService() AcademicDegreeService {
	return &academicdegreeservice{
		db: config.DB,
	}
}

func (s *academicdegreeservice) GetAcademicDegreeByAcademicID(
	ctx context.Context,
	id int,
) ([]response.AcademicDegreeResponse, error) {

	if id <= 0 {
		return nil, apperror.New(
			apperror.CodeInvalidInput,
			"academic id is required",
			nil,
		)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.AcademicDegreeResponse

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("academic_degrees ad").
			Joins("LEFT JOIN majors m ON m.id = ad.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
			Joins("LEFT JOIN academics a ON a.id = ad.academic_id")
	}

	dataQuery := base().
		Select(`
			ad.id AS id,
			ad.uuid AS uuid,
			a.id AS academic_id,
			a.code AS academic_code,
			a.name AS academic_name,
			f.id AS faculty_id,
			d.id AS department_id,
			m.id AS major_id,
			m.code AS major_code,
			m.name AS major_name,
			p.id AS programme_id,
			p.name AS programme_name,
			ad.name AS name,
			ad.monthly_fee AS monthly_fee,
			ad.quarterly_fee AS quarterly_fee,
			ad.semesterly_fee AS semesterly_fee,
			ad.yearly_fee AS yearly_fee,
			ad.description AS description,
			ad.active AS active
		`).
		Where("ad.academic_id = ?", id)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
func (s *academicdegreeservice) GetAcademicDegree(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AcademicDegreeResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.AcademicDegreeResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("academic_degrees ad").
			Joins("LEFT JOIN majors m ON m.id = ad.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
			Joins("LEFT JOIN academics a ON a.id = ad.academic_id")
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
		ad.id AS id,
		ad.uuid AS uuid,
		a.id AS academic_id,
		a.code AS academic_code,
		a.name AS academic_name,
		f.id AS faculty_id,
		d.id AS department_id,
		m.id AS major_id,
		m.code AS major_code,
		m.name AS major_name,
		p.id AS programme_id,
		p.name AS programme_name,
		ad.name AS name,
		ad.monthly_fee AS monthly_fee,
		ad.quarterly_fee AS quarterly_fee,
		ad.semesterly_fee AS semesterly_fee,
		ad.yearly_fee AS yearly_fee,
		ad.description AS description,
		ad.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *academicdegreeservice) CreateAcademicDegree(ctx context.Context, input request.AcademicDegreeRequestCreate) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if err := validate.Struct(input); err != nil {
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newData := model.AcademicDegree{
		UUIDBase:      base.UUIDBase{UUID: helper.GenerateUUID()},
		AcademicID:    input.AcademicID,
		MajorID:       input.MajorID,
		Name:          input.Name,
		MonthlyFee:    input.MonthlyFee,
		QuarterlyFee:  input.QuarterlyFee,
		SemesterlyFee: input.SemesterlyFee,
		YearlyFee:     input.YearlyFee,
		Description:   input.Description,
		Active:        true,
	}

	if err := s.db.WithContext(ctx).Create(&newData).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}
	return nil
}

func (s *academicdegreeservice) UpdateAcademicDegree(ctx context.Context, id string, input request.AcademicDegreeRequestUpdate) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperror.New(apperror.CodeInvalidInput, "academic degree id is required", nil)
	}

	if err := validate.Struct(input); err != nil {
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	updates := map[string]interface{}{}

	if input.AcademicID != nil {
		updates["academic_id"] = *input.AcademicID
	}

	if input.MajorID != nil {
		updates["major_id"] = *input.MajorID
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "academic degree name cannot be empty", nil)
		}
		updates["name"] = name
	}

	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}

	if input.MonthlyFee != nil {
		updates["monthly_fee"] = *input.MonthlyFee
	}

	if input.QuarterlyFee != nil {
		updates["quarterly_fee"] = *input.QuarterlyFee
	}

	if input.SemesterlyFee != nil {
		updates["semesterly_fee"] = *input.SemesterlyFee
	}

	if input.YearlyFee != nil {
		updates["yearly_fee"] = *input.YearlyFee
	}

	if len(updates) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "no fields provided to update", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).
		Model(&model.AcademicDegree{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "update")
	}

	if result.RowsAffected == 0 {
		return apperror.NotFound("academic degree not found", nil)
	}

	return nil
}

func (s *academicdegreeservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.AcademicDegree](ctx, s.db, id)
}
