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
	"strings"

	"gorm.io/gorm"
)

type ClassCurriculumnService interface {
	CreateClassCurriculumn(ctx context.Context, input request.ClassCurriculumnRequestCreate) error
}

type classcurriculmnservice struct {
	db *gorm.DB
}

func NewClassCurriculumnService() ClassCurriculumnService {
	return &classcurriculmnservice{
		db: config.DB,
	}
}

func (s *classcurriculmnservice) CreateClassCurriculumn(ctx context.Context, input request.ClassCurriculumnRequestCreate) error {
	name := strings.TrimSpace(input.Name)
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newclasscurriculumn := model.ClassCurriculumn{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		Name:    name,
		MajorID: input.MajorID,
		TermID:  input.TermID,
		Active:  true,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newclasscurriculumn).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create class curriculmn", nil)
		}

		if len(input.ClassCurriculumnDetailRequestCreate) > 0 {
			detail := make([]model.ClassCurriculumnDetail, 0, len(input.ClassCurriculumnDetailRequestCreate))
			for _, c := range input.ClassCurriculumnDetailRequestCreate {
				detail = append(detail, model.ClassCurriculumnDetail{
					UUIDBase: base.UUIDBase{
						UUID: helper.GenerateUUID(),
					},
					ClassCurriculumnID: newclasscurriculumn.ID,
					SemesterID:         c.SemesterID,
					StudyYearID:        c.StudyYearID,
					AcademicShiftID:    c.AcademicShiftID,
					MidtermDate:        c.MidtermDate,
					FinalDate:          c.FinalDate,
					TotalStudent:       nil,
					TypeClass:          c.TypeClass,
				})
			}
			if err := tx.Create(&detail).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create dedtail", nil)
			}
		}
		return nil
	})
	return err
}
