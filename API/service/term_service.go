package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TermService interface {
	GetTerm(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TermResponse, *model.PaginationMetadata, error)
	CreateTerm(ctx context.Context, input request.TermRequestCreate) error
	UpdateTerm(ctx context.Context, id string, input request.TermRequestUpdate) error
	Toggle(ctx context.Context, id string) error
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
		t.uuid AS uuid,
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
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation code is required", nil)
	}

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation name is required", nil)
	}

	if input.GenerationID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "generation id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	nextIndex := 1

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lastTerm model.Term
		err := tx.
			Where("generation_id = ?", input.GenerationID).
			Order("id DESC").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastTerm).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			nextIndex = 1
		case err != nil:
			return helper.MapAcademicError(err, "LOOKUP_LAST_GENERATION")
		default:
			nextIndex = lastTerm.Index + 1
		}

		newdata := model.Term{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			GenerationID: input.GenerationID,
			Code:         code,
			Name:         name,
			Index:        nextIndex,
			StartDate:    input.StartDate,
			EndDate:      nil,
			Description:  input.Description,
			Active:       true,
		}

		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapAcademicError(err, "CREATE")
		}
		return nil
	})
	return err
}

func (s *termservice) UpdateTerm(ctx context.Context, id string, input request.TermRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.GenerationID != nil {
		updates["generation_id"] = *input.GenerationID
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
		Model(&model.Term{}).
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

func (s *termservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Term](ctx, s.db, id)
}
