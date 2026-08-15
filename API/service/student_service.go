package service

import (
	"context"
	"strings"

	"mysql/config"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/request"
	"mysql/utils"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var studentValidator = validator.New()

type StudentService interface {
	CreateStudentService(ctx context.Context, input request.StudentRequestCreate) error
}

type studentService struct {
	db *gorm.DB
}

func NewStudentService() StudentService {
	return &studentService{
		db: config.DB,
	}
}

func (s *studentService) CreateStudentService(ctx context.Context, input request.StudentRequestCreate) error {
	input.NameKh = strings.TrimSpace(input.NameKh)
	input.NameEn = strings.TrimSpace(input.NameEn)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Nationality = strings.TrimSpace(input.Nationality)

	if err := studentValidator.Struct(input); err != nil {
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	student := model.Student{
		GroupID:          input.GroupID,
		NameKh:           input.NameKh,
		NameEn:           input.NameEn,
		DateOfBirth:      input.DateOfBirth,
		Gender:           input.Gender,
		Nationality:      input.Nationality,
		Phone:            input.Phone,
		VillageID:        input.VillageID,
		Occupation:       input.Occupation,
		AcademicStreamID: input.AcademicStreamID,
	}

	family := model.StudentFamily{
		FatherName:        input.FatherName,
		FatherEnglishName: input.FatherEnglishName,
		FatherAge:         input.FatherAge,
		FatherIsAlive:     input.FatherIsAlive,
		FatherPhoneNumber: input.FatherPhoneNumber,
		FatherOccupation:  input.FatherOccupation,
		FatherWorkplace:   input.FatherWorkplace,
		MotherName:        input.MotherName,
		MotherEnglishName: input.MotherEnglishName,
		MotherAge:         input.MotherAge,
		MotherIsAlive:     input.MotherIsAlive,
		MotherPhoneNumber: input.MotherPhoneNumber,
		MotherOccupation:  input.MotherOccupation,
		MotherWorkplace:   input.MotherWorkplace,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&student).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student", nil)
		}

		family.StudentID = student.ID
		if err := tx.Create(&family).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student family", nil)
		}

		if len(input.StudentEducationRequestCreate) > 0 {
			educations := make([]model.StudentEducation, 0, len(input.StudentEducationRequestCreate))
			for _, e := range input.StudentEducationRequestCreate {
				educations = append(educations, model.StudentEducation{
					StudentID:       student.ID,
					Level:           e.Level,
					SchoolName:      e.SchoolName,
					VillageID:       e.VillageID,
					StartDate:       e.StartDate,
					EndDate:         e.EndDate,
					CertificateDate: e.CertificateDate,
					Score:           e.Score,
					Gpa:             e.Gpa,
					Grade:           e.Grade,
				})
			}
			if err := tx.Create(&educations).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create student educations", nil)
			}
		}

		if len(input.StudentDocumentRequestCreate) > 0 {
			documents := make([]model.StudentDocument, 0, len(input.StudentDocumentRequestCreate))
			for _, d := range input.StudentDocumentRequestCreate {
				documents = append(documents, model.StudentDocument{
					StudentID:      student.ID,
					DocumentTypeID: d.DocumentTypeID,
					RequiredQty:    d.RequiredQty,
					RecieveQty:     d.RecieveQty,
					Remark:         d.Remark,
				})
			}
			if err := tx.Create(&documents).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create student document", nil)
			}
		}

		return nil
	})

	return err
}
