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

	"gorm.io/gorm"
)

type MajorTermService interface {
	CreateMajorTerm(ctx context.Context, input request.MajorTermReqeustCreate) error
	GetMajorTerm(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.MajorTermResponse, *model.PaginationMetadata, error)
	UpdateMajorTerm(ctx context.Context, id string, input request.MajorTermReqeustUpdate) error
	Toggle(ctx context.Context, id string) error
}

type majortermservice struct {
	db *gorm.DB
}

func NewMajorTermService() MajorTermService {
	return &majortermservice{
		db: config.DB,
	}
}

func (s *majortermservice) CreateMajorTerm(ctx context.Context, input request.MajorTermReqeustCreate) error {
	if input.TermID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "term is required", nil)
	}
	if len(input.MajorID) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "major_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := make([]model.MajorTerm, 0, len(input.MajorID))
	for _, majorID := range input.MajorID {
		newdata = append(newdata, model.MajorTerm{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			MajorID: majorID,
			TermID:  input.TermID,
			Active:  true,
		})
	}

	if err := s.db.WithContext(ctx).Create(&newdata).Error; err != nil {
		return apperror.New(apperror.CodeInternal, "failed to create major term", err)
	}

	return nil
}

func (s *majortermservice) GetMajorTerm(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.MajorTermResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.MajorTermResponse
	var total int64

	baseQuery := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("major_terms mt").
			Joins("LEFT JOIN majors m ON m.id = mt.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
			Joins("LEFT JOIN terms t ON t.id = mt.term_id").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id")
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
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("g.id = ?", v)
		}
		if v, ok := filter["term_id"]; ok && v != "" {
			tx = tx.Where("t.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(baseQuery()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count major terms: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	// NOTE: id/uuid come from mt (the major_terms junction row), not from the
	// major itself - UpdateMajorTerm looks records up by major_terms.uuid, so
	// the identity returned here has to match that table or edits will 404.
	dataQuery := applyFilters(baseQuery()).Select(`
		mt.id AS id,
		mt.uuid AS uuid,
		m.id AS major_id,
		m.code AS code,
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
		p.name AS programme_name,
		t.id AS term_id,
		t.code AS term_code,
		t.name AS term_name,
		g.id AS generation_id,
		g.code AS generation_code,
		g.name AS generation_name,
		a.id AS academic_id,
		a.code AS academic_code,
		a.name AS academic_name
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch major terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *majortermservice) UpdateMajorTerm(ctx context.Context, id string, input request.MajorTermReqeustUpdate) error {
	if id == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}
	if input.TermID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "term is required", nil)
	}
	if input.MajorID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "major_id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var existing model.MajorTerm
	if err := s.db.WithContext(ctx).Where("uuid = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.New(apperror.CodeNotFound, "major term not found", nil)
		}
		return apperror.New(apperror.CodeInternal, "failed to fetch major term", err)
	}

	existing.MajorID = input.MajorID
	existing.TermID = input.TermID

	if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return apperror.New(apperror.CodeInternal, "failed to update major term", err)
	}

	return nil
}

func (s *majortermservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.MajorTerm](ctx, s.db, id)
}
