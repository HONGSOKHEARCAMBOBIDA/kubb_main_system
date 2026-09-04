package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/utils"

	"gorm.io/gorm"
)

type ClassRepresentativeService interface {
	CreateClassRepresentative(ctx context.Context, input request.ClassRepresentativeRequestCreate) error
}

type classrepresentativeservice struct {
	db *gorm.DB
}

func NewClassRepresentative() ClassRepresentativeService {
	return &classrepresentativeservice{
		db: config.DB,
	}
}

func (s *classrepresentativeservice) CreateClassRepresentative(ctx context.Context, input request.ClassRepresentativeRequestCreate) error {
	if input.ClassCurriculumnDetailID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "classcurriculumn detail id is required", nil)
	}

	if len(input.ClassRepresentativeRequest) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "class representative is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		classrepresentatives := make([]model.ClassRepresentative, 0, len(input.ClassRepresentativeRequest))
		for _, c := range input.ClassRepresentativeRequest {
			classrepresentatives = append(classrepresentatives, model.ClassRepresentative{
				UUIDBase: base.UUIDBase{
					UUID: helper.GenerateUUID(),
				},
				ClassCurriculumnDetailID: input.ClassCurriculumnDetailID,
				StudentTermID:            c.StudentTermID,
				StartDate:                c.StartDate,
				EndDate:                  c.EndDate,
				Isactive:                 true,
			})
		}

		if err := tx.Create(&classrepresentatives).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create", err)
		}
		return nil
	})

	return err
}
