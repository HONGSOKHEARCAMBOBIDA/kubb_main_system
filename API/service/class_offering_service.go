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

type ClassOfferingService interface {
	CreateClassOffering(ctx context.Context, input request.ClassOfferingRequestCreate) error
}

type classofferingservice struct {
	db *gorm.DB
}

func NewClassOfferingService() ClassOfferingService {
	return &classofferingservice{
		db: config.DB,
	}
}

func (s *classofferingservice) CreateClassOffering(ctx context.Context, input request.ClassOfferingRequestCreate) error {
	if input.ClassCurriculumnDetailID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "academic_id is required", nil)
	}
	if len(input.ClassOfferingRequest) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "class_offering is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		classoffering := make([]model.ClassOffering, 0, len(input.ClassOfferingRequest))
		for _, c := range input.ClassOfferingRequest {
			if c.SubjectID == 0 {
				return apperror.New(apperror.CodeInvalidInput, "subject_id is required", nil)
			}
			classoffering = append(classoffering, model.ClassOffering{
				UUIDBase: base.UUIDBase{
					UUID: helper.GenerateUUID(),
				},
				ClassCurriculumnDetailID:  input.ClassCurriculumnDetailID,
				SubjectID:                 c.SubjectID,
				Credit:                    c.Credit,
				PassingScore:              c.PassingScore,
				TotalHour:                 c.TotalHour,
				Status:                    model.StatusClassOfferingOpen,
				TotalAttendanceForRexam:   c.TotalAttendanceForRexam,
				TotalAttendanceForRelearn: c.TotalAttendanceForRelearn,
				Description:               c.Description,
			})
		}

		if err := tx.Create(&classoffering).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create", err)
		}
		return nil
	})

	return err
}
