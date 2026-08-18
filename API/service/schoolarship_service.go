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

var validateschoolarship = validator.New()

type SchoolarshipService interface {
	GetSchoolarship(
		ctx context.Context,
		pf request.Pagination,
		filter map[string]string,
	) ([]response.SchoolarshipResponse, *model.PaginationMetadata, error)

	CreateSchoolarship(
		ctx context.Context,
		input request.SchoolarshipRequestCreate,
	) error

	UpdateSchoolarship(
		ctx context.Context,
		id string,
		input request.SchoolarshipRequestUpdate,
	) error

	Toggle(ctx context.Context, id string) error
}

type schoolarshipservice struct {
	db *gorm.DB
}

func NewSchoolarshipService() SchoolarshipService {
	return &schoolarshipservice{
		db: config.DB,
	}
}

func (s *schoolarshipservice) GetSchoolarship(
	ctx context.Context,
	pf request.Pagination,
	filter map[string]string,
) ([]response.SchoolarshipResponse, *model.PaginationMetadata, error) {

	helper.NormalizePagination(&pf)

	var data []response.SchoolarshipResponse
	var total int64

	baseQuery := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("scholarships s")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("s.name LIKE ?", "%"+v+"%")
		}

		if v, ok := filter["code"]; ok && v != "" {
			tx = tx.Where("s.code LIKE ?", "%"+v+"%")
		}

		return tx
	}

	// Count
	if err := applyFilters(baseQuery()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count scholarships: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(baseQuery()).
		Select(`
			s.id AS id,
			s.uuid AS uuid,
			s.code AS code,
			s.name AS name,
			s.discount_type AS discount_type,
			s.discount_percentage AS discount_percentage,
			s.discount_amount AS discount_amount,
			s.description AS description,
			s.active AS active
		`)

	if err := dataQuery.
		Offset(offset).
		Limit(pf.PageSize).
		Scan(&data).
		Error; err != nil {
		return nil, nil, fmt.Errorf("fetch scholarships: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *schoolarshipservice) CreateSchoolarship(
	ctx context.Context,
	input request.SchoolarshipRequestCreate,
) error {

	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if err := validateschoolarship.Struct(input); err != nil {
		return apperror.New(
			apperror.CodeInvalidInput,
			err.Error(),
			nil,
		)
	}

	discountType := input.DiscountType

	if discountType == "" {
		discountType = model.DiscountPercentage
	}

	switch discountType {
	case model.DiscountAmount, model.DiscountPercentage:
	default:
		return apperror.New(
			apperror.CodeInvalidInput,
			"invalid discount_type",
			nil,
		)
	}

	// Prevent inconsistent discount values
	if discountType == model.DiscountPercentage &&
		input.DiscountAmount != 0 {

		return apperror.New(
			apperror.CodeInvalidInput,
			"discount_amount must be 0 when discount_type is percentage",
			nil,
		)
	}

	if discountType == model.DiscountAmount &&
		input.DiscountPercentage != 0 {

		return apperror.New(
			apperror.CodeInvalidInput,
			"discount_percentage must be 0 when discount_type is amount",
			nil,
		)
	}

	ctx, cancel := context.WithTimeout(
		ctx,
		utils.DefaultQueryTimeout,
	)
	defer cancel()

	newData := model.Schoolarship{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		Code:               input.Code,
		Name:               input.Name,
		DiscountType:       discountType,
		DiscountPercentage: input.DiscountPercentage,
		DiscountAmount:     input.DiscountAmount,
		Description:        input.Description,
		Active:             true,
	}

	if err := s.db.
		WithContext(ctx).
		Create(&newData).
		Error; err != nil {

		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}

func (s *schoolarshipservice) UpdateSchoolarship(
	ctx context.Context,
	id string,
	input request.SchoolarshipRequestUpdate,
) error {

	id = strings.TrimSpace(id)

	if id == "" {
		return apperror.New(
			apperror.CodeInvalidInput,
			"id is required",
			nil,
		)
	}

	if input.Code != nil {
		trimmed := strings.TrimSpace(*input.Code)
		input.Code = &trimmed
	}

	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		input.Name = &trimmed
	}

	if input.Description != nil {
		trimmed := strings.TrimSpace(*input.Description)
		input.Description = &trimmed
	}

	if err := validateschoolarship.Struct(input); err != nil {
		return apperror.New(
			apperror.CodeInvalidInput,
			err.Error(),
			nil,
		)
	}

	updates := map[string]interface{}{}

	if input.Code != nil {
		if *input.Code == "" {
			return apperror.New(
				apperror.CodeInvalidInput,
				"code is required",
				nil,
			)
		}

		updates["code"] = *input.Code
	}

	if input.Name != nil {
		if *input.Name == "" {
			return apperror.New(
				apperror.CodeInvalidInput,
				"name is required",
				nil,
			)
		}

		updates["name"] = *input.Name
	}

	if input.DiscountType != nil {

		switch *input.DiscountType {

		case model.DiscountPercentage:

			if input.DiscountAmount != nil &&
				*input.DiscountAmount != 0 {

				return apperror.New(
					apperror.CodeInvalidInput,
					"discount_amount must not be set when discount_type is percentage",
					nil,
				)
			}

			updates["discount_type"] = *input.DiscountType

		case model.DiscountAmount:

			if input.DiscountPercentage != nil &&
				*input.DiscountPercentage != 0 {

				return apperror.New(
					apperror.CodeInvalidInput,
					"discount_percentage must not be set when discount_type is amount",
					nil,
				)
			}

			updates["discount_type"] = *input.DiscountType

		default:

			return apperror.New(
				apperror.CodeInvalidInput,
				"invalid discount_type",
				nil,
			)
		}
	}

	if input.DiscountPercentage != nil {
		updates["discount_percentage"] = *input.DiscountPercentage
	}

	if input.DiscountAmount != nil {
		updates["discount_amount"] = *input.DiscountAmount
	}

	if input.Description != nil {
		updates["description"] = *input.Description
	}

	if len(updates) == 0 {
		return apperror.New(
			apperror.CodeInvalidInput,
			"no fields provided to update",
			nil,
		)
	}

	ctx, cancel := context.WithTimeout(
		ctx,
		utils.DefaultQueryTimeout,
	)
	defer cancel()

	result := s.db.
		WithContext(ctx).
		Model(&model.Schoolarship{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "UPDATE")
	}

	if result.RowsAffected == 0 {
		return apperror.New(
			apperror.CodeNotFound,
			"schoolarship not found",
			nil,
		)
	}

	return nil
}

func (s *schoolarshipservice) Toggle(
	ctx context.Context,
	id string,
) error {

	return utils.ToggleStatus[model.Schoolarship](
		ctx,
		s.db,
		id,
	)
}
