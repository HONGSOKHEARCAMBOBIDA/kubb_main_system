package service

import (
	"context"
	"fmt"

	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type TermService interface {
	GetTerm(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TermResponse, *model.PaginationMetadata, error)
	CreateTerm(ctx context.Context, input request.TermRequestCreate) error
	UpdateTerm(ctx context.Context, id string, input request.TermRequestUpdate) error
	Toggle(ctx context.Context, id string)
}

type termservice struct {
	db *gorm.DB
}

func NewTermService() TermService {
	return &termservice{
		db: config.DB,
	}
}

func (s *termservice) GetTerm(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TermResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.TermResponse
	var total int64

	// Base query WITH joins applied consistently to both count and data queries.
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("terms t").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("g.id = ?", v)
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
		t.id AS id,
		g.id AS generation_id,
		g.code AS generation_code,
		g.name AS generation_name,
		a.id AS academic_id,
		a.code AS academic_code,
		a.name AS academic_name,
		t.code AS code,
		t.name AS name,
		t.index AS ` + "`index`" + `,
		t.start_date AS start_date,
		t.end_date AS end_date,
		t.description AS description,
		t.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *termservice) CreateTerm(ctx context.Context, input request.TermRequestCreate) error {

}
