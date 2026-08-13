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

var validatefeediscountgroup = validator.New()

type FeeDiscountGroupService interface {
	GetFeeDiscountGroup(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.FeeDiscountGroupResponse, *model.PaginationMetadata, error)
	CreateFeeDiscountGroup(ctx context.Context, input request.FeeDiscountGroupRequestCreate) error
	UpdateFeeDiscountGroup(ctx context.Context, id string, input request.FeeDiscountGroupRequestUpdate) error
	Toggle(ctx context.Context, id string) error
}

type feediscountgroupservice struct {
	db *gorm.DB
}

func NewFeeDiscountGroupService() FeeDiscountGroupService {
	return &feediscountgroupservice{
		db: config.DB,
	}
}

func (s *feediscountgroupservice) GetFeeDiscountGroup(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.FeeDiscountGroupResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.FeeDiscountGroupResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("fee_discount_groups f")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("f.name LIKE ?", v)
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
		f.code AS code,
		f.name AS name,
		f.discount_type AS discount_type,
		f.discount_percentage AS discount_percentage,
		f.discount_amount AS discount_amount,
		f.description AS description,
		f.active AS active
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil

}

func (s *feediscountgroupservice) CreateFeeDiscountGroup(ctx context.Context, input request.FeeDiscountGroupRequestCreate) error {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	if err := validatefeediscountgroup.Struct(input); err != nil {
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	discounttype := input.DiscountType
	if discounttype == "" {
		discounttype = model.DiscountPercentage
	}

	switch discounttype {
	case model.DiscountAmount, model.DiscountPercentage:
	default:
		return apperror.New(apperror.CodeInvalidInput, "invalid ", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newData := model.FeeDiscountGroup{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		Code:               input.Code,
		Name:               input.Name,
		DiscountType:       discounttype,
		DiscountPercentage: input.DiscountPercentage,
		DiscountAmount:     input.DiscountAmount,
		Description:        input.Description,
		Active:             true,
	}
	if err := s.db.WithContext(ctx).Create(&newData).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}
