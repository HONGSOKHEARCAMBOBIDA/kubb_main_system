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
		return apperror.New(apperror.CodeInvalidInput, "class offering is required", nil)
	}

	if len(input.CourseRegistrationRequest) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "studentterm is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		courseregister := make([]model.CourseRegistrations, 0, len(input.CourseRegistrationRequest))
		for _, c := range input.CourseRegistrationRequest {
			if c.StudentTermID == 0 {
				return apperror.New(apperror.CodeInvalidInput, "student term is required", nil)
			}
			courseregister = append(courseregister, model.CourseRegistrations{
				UUIDBase: base.UUIDBase{
					UUID: helper.GenerateUUID(),
				},
				StudentTermID:    c.StudentTermID,
				ClassOfferingID:  input.ClassOfferingID,
				RegistrationDate: time.Now().Format("2006-01-02"),
				Status:           "PENDING",
			})
		}
		if err := tx.Create(&courseregister).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create", err)
		}

		studentgrades := make([]model.StudentGrade, 0, len(courseregister))
		for _, cr := range courseregister {
			studentgrades = append(studentgrades, model.StudentGrade{
				UUIDBase: base.UUIDBase{
					UUID: helper.GenerateUUID(),
				},
				CourseRegistrationID: cr.ID,
				TotalScore:           nil,
				LetterGrade:          nil,
				GradePoint:           nil,
				Status:               "PENDING",
			})
		}
		if err := tx.Create(&studentgrades).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create", err)
		}

		var classoffer model.ClassOffering
		if err := tx.First(&classoffer, input.ClassOfferingID).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to fetch class offering", err)
		}

		var gradecomponent []model.GradeComponent
		if err := tx.Where("subject_id = ?", classoffer.SubjectID).Find(&gradecomponent).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to fetch grade components", err)
		}

		fmt.Println("classoffer.SubjectID:", classoffer.SubjectID, "gradecomponent count:", len(gradecomponent))

		if len(gradecomponent) > 0 {
			studentgradedetails := make([]model.StudentGradeDetail, 0, len(studentgrades)*len(gradecomponent))
			for _, sg := range studentgrades {
				for _, gc := range gradecomponent {
					studentgradedetails = append(studentgradedetails, model.StudentGradeDetail{
						UUIDBase: base.UUIDBase{
							UUID: helper.GenerateUUID(),
						},
						StudentGradeID:   sg.ID,
						GradeComponentID: gc.ID,
						Score:            nil,
					})
				}
			}
			if err := tx.Create(&studentgradedetails).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create", err)
			}
		}

		return nil
	})
	return err
}
