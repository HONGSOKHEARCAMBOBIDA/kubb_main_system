package service

import (
	"context"
	"fmt"
	"strings"

	"mysql/config"
	"mysql/constant/apperror"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var studentValidator = validator.New()

type StudentService interface {
	CreateStudent(ctx context.Context, input request.StudentRequestCreate) error
	GetStudent(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.StudentResponse, *model.PaginationMetadata, error)
}

type studentService struct {
	db *gorm.DB
}

func NewStudentService() StudentService {
	return &studentService{
		db: config.DB,
	}
}

func (s *studentService) CreateStudent(ctx context.Context, input request.StudentRequestCreate) error {
	input.NameKh = strings.TrimSpace(input.NameKh)
	input.NameEn = strings.TrimSpace(input.NameEn)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Nationality = strings.TrimSpace(input.Nationality)
	username := strings.ToLower(input.NameEn)
	email := helper.GenerateEmail(username, 168)
	if err := studentValidator.Struct(input); err != nil {
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	student := model.Student{
		GroupID:          input.GroupID,
		UserName:         username,
		Email:            email,
		Password:         utils.HasPassword("KUBB"),
		NameKh:           input.NameKh,
		NameEn:           input.NameEn,
		DateOfBirth:      input.DateOfBirth,
		Gender:           input.Gender,
		Nationality:      input.Nationality,
		Phone:            input.Phone,
		Status:           share.Created,
		VillageID:        input.VillageID,
		Occupation:       input.Occupation,
		AcademicStreamID: input.AcademicStreamID,
		TelegramUsername: nil,
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

		student.Code = helper.GenerateCode("STU", uint(student.ID))

		if err := tx.Save(&student).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "faild to update code", nil)
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

func (s *studentService) GetStudent(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.StudentResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.StudentResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("students s").
			Joins("LEFT JOIN fee_discount_groups f ON f.id = s.group_id").
			Joins("LEFT JOIN academic_streams asd ON asd.id = s.academic_stream_id").
			Joins("LEFT JOIN villages sv ON sv.id = s.village_id").
			Joins("LEFT JOIN communes sc ON sc.id = sv.communce_id").
			Joins("LEFT JOIN districts sd ON sd.id = sc.district_id").
			Joins("LEFT JOIN provinces sp ON sp.id = sd.province_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("s.name_kh LIKE ?", "%"+v+"%")
		}

		return tx
	}

	if err := applyFilters(base()).
		Count(&total).Error; err != nil {

		return nil, nil, fmt.Errorf("count terms: %w", err)
	}	

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).
	Select(`
		s.id AS id,
		s.uuid AS uuid,
		f.id AS group_id,
		f.code AS group_code,
		f.name AS group_name,
		f.discount_type AS discount_type,
		f.discount_percentage AS discount_percentage,
		f.discount_amount AS discount_amount,
		s.code AS code,
		s.username AS username,
		s.email AS email,
		s.name_kh AS name_kh,
		s.name_en AS name_en,
		s.date_of_birth AS date_of_birth,
		s.gender AS gender,
		s.nationality AS nationality,
		s.phone AS phone,
		s.status AS status,
		sv.id AS village_id,
		sv.name_kh AS villlage_name_kh,
		sc.id AS communce_id,
		sc.name_kh AS communce_name,
		sd.id AS district_id,
		sd.name_kh AS distirct_name,
		sp.id AS province_id,
		sp.name_kh AS province_name,
		s.occupation AS occupation,
		asd.id AS academic_stream_id,
		asd.name AS academic_stream_name
	`).Order("s.id DESC").Offset(offset).Limit(pf.PageSize)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch student: %w", err)
	}

	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	studentIDs := make([]int,0,len(data))

	for _, student := range data {
		studentIDs = append(studentIDs,student.ID)
	}

	var studentfamily []model.StudentFamily

	if err := s.db.WithContext(ctx).
	Table("student_families sf").
	Joins("INNER JOIN students s ON s.id = sf.student_id").Where("sf.student_id IN ?",studentIDs).
	Select(`
		sf.id AS id,
		sf.student_id AS student_id,
		sf.father_name AS father_name,
		sf.father_english_name AS father_english_name,
		sf.father_age AS father_age,
		sf.father_is_alive AS father_is_alive,
		sf.father_phone_number AS father_phone_number,
		sf.father_occupation AS father_occupation,
		sf.father_workplace AS father_workplace,
		sf.mother_name AS mother_name,
		sf.mother_english_name AS mother_english_name,
		sf.mother_age AS mother_age,
		sf.mother_is_alive AS mother_is_alive,
		sf.mother_phone_number AS mother_phone_number,
		sf.mother_phone_number AS mother_phone_number,
		sf.mother_occupation AS mother_occupation,
		sf.mother_workplace AS mother_workplace
	`).Scan(&studentfamily).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch term student family: %w", err)
	}

	var studenteducation []response.StudentEducationResponse

	if err := s.db.WithContext(ctx).
	Table()

}
