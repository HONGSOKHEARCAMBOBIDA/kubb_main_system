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

type AttendanceService interface {
	CreateAttendance(ctx context.Context, input request.AttendanceRequestCreate) error
}

type attendanceservice struct {
	db *gorm.DB
}

func NewAttendanceService() AttendanceService {
	return &attendanceservice{
		db: config.DB,
	}
}

func (s *attendanceservice) CreateAttendance(ctx context.Context, input request.AttendanceRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Attendance{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		ScheduleID:     input.ScheduleID,
		AttendanceDate: input.AttendanceDate,
		Topic:          input.Topic,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create attendance", nil)
		}
		if len(input.AttendanceDetailRequestCreate) > 0 {
			detail := make([]model.AttendanceDetail, 0, len(input.AttendanceDetailRequestCreate))
			for _, d := range input.AttendanceDetailRequestCreate {
				detail = append(detail, model.AttendanceDetail{
					UUIDBase: base.UUIDBase{
						UUID: helper.GenerateUUID(),
					},
					AttendanceID:          newdata.ID,
					CourseRegistrationsID: d.CourseRegistrationsID,
					Status:                d.Status,
					Note:                  d.Note,
				})
			}
			if err := tx.Create(&detail).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create attendance detail", nil)
			}
		}
		return nil
	})
	return err
}
