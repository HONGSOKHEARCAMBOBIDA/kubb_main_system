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
	"time"

	"gorm.io/gorm"
)

type CourseRegistrationService interface {
	CreateCourseRegistration(ctx context.Context, input request.CourseRegistrationRequestCreate) error
}

type couseregistrationservice struct {
	db *gorm.DB
}

func NewCourseRegistrationService() CourseRegistrationService {
	return &couseregistrationservice{
		db: config.DB,
	}
}

func (s *couseregistrationservice) CreateCourseRegistration(ctx context.Context, input request.CourseRegistrationRequestCreate) error {
	if input.ClassOfferingID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "class offerig is required", nil)
	}

	if len(input.CourseRegistrationRequest) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "studentterm is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		coursregister := make([]model.CourseRegistrations, 0, len(input.CourseRegistrationRequest))
		for _, c := range input.CourseRegistrationRequest {
			if c.StudentTermID == 0 {
				return apperror.New(apperror.CodeInvalidInput, "student term is required", nil)
			}
			coursregister = append(coursregister, model.CourseRegistrations{
				UUIDBase: base.UUIDBase{
					UUID: helper.GenerateUUID(),
				},
				StudentTermID:    c.StudentTermID,
				ClassOfferingID:  input.ClassOfferingID,
				RegistrationDate: time.Now().Format("2006-01-02"),
				Status:           "PENDING",
			})
		}
		if err := tx.Create(&coursregister).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create", err)
		}
		return nil
	})
	return err
}
