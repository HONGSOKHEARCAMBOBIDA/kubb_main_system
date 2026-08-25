package service

import (
	"context"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strings"

	"gorm.io/gorm"
)

type TeacherService interface {
	CreateTeacher(ctx context.Context, input request.TeacherRequestCreate) error
	GetTeacher(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TeacherResponse, error)
	UpdateTeacher(ctx context.Context, uuid string, input request.TeacherRequestUpdate) error
	Toggle(ctx context.Context, uuid string) error
}

type teacherservice struct {
	db *gorm.DB
}

func NewTeacherService() TeacherService {
	return &teacherservice{
		db: config.DB,
	}
}

func (s *teacherservice) CreateTeacher(ctx context.Context, input request.TeacherRequestCreate) error {
	email := strings.TrimSpace(input.Email)
	name := strings.TrimSpace(input.Name)
	if input.FacultyID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "faculty id is required", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newteacher := model.Teacher{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		Email:       email,
		Password:    utils.HasPassword("KUBB"),
		Name:        name,
		Dob:         input.Dob,
		Pob:         input.Pob,
		Gender:      input.Gender,
		Nationality: input.Nationality,
		Address:     input.Address,
		Phone:       input.Phone,
		FacultyID:   input.FacultyID,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newteacher).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create teacher", nil)
		}
		newteacher.Code = helper.GenerateCode("T", uint(newteacher.ID))
		return nil
	})
	return err
}

func (s *teacherservice) GetTeacher(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TeacherResponse, error) {
	helper.NormalizePagination(&pf)
	var data []response.TeacherResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("teachers t").
			Joins("LEFT JOIN faculties f ON f.id = t.faculty_id")
	}

}
