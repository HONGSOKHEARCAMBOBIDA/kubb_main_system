package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

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
	UpdateStudent(ctx context.Context, studentID int, input request.StudentRequestUpdate, id int) error
	GetCourseRegistration(ctx context.Context, filter map[string]string) ([]response.CourseRegistrationResponse, error)
	GetStudentCategory(ctx context.Context) ([]model.StudentCategory, error)
}

type studentService struct {
	db *gorm.DB
}

func NewStudentService() StudentService {
	return &studentService{
		db: config.DB,
	}
}

func (s *studentService) GetStudentCategory(ctx context.Context) ([]model.StudentCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.StudentCategory
	err := s.db.WithContext(ctx).Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch filingcabinet", err)
	}
	return data, nil
}

func (s *studentService) UpdateStudent(ctx context.Context, studentID int, input request.StudentRequestUpdate, id int) error {
	input.NameKh = strings.TrimSpace(input.NameKh)
	input.NameEn = strings.TrimSpace(input.NameEn)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Nationality = strings.TrimSpace(input.Nationality)
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return err
	}

	newStatus := fmt.Sprintf("Update by %s", user.NameKh)

	if err := studentValidator.Struct(input); err != nil {
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var student model.Student
		if err := tx.First(&student, studentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "student not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch student", nil)
		}

		student.GroupID = input.GroupID
		student.StudentCategoryID = input.StudentCategoryID
		student.NameKh = input.NameKh
		student.NameEn = input.NameEn
		student.DateOfBirth = input.DateOfBirth
		student.Gender = input.Gender
		student.Nationality = input.Nationality
		student.Phone = input.Phone
		student.VillageID = &input.VillageID
		student.Occupation = input.Occupation
		student.AcademicStreamID = input.AcademicStreamID
		student.Status = newStatus
		if err := tx.Save(&student).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}

		if err := tx.Where("student_id = ?", student.ID).
			Delete(&model.StudentFamily{}).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to clear student educations", nil)
		}

		if len(input.StudentFamilyRequestUpdate) > 0 {
			family := make([]model.StudentFamily, 0, len(input.StudentFamilyRequestUpdate))
			for _, f := range input.StudentFamilyRequestUpdate {
				family = append(family, model.StudentFamily{
					StudentID:         student.ID,
					FatherName:        f.FatherName,
					FatherEnglishName: f.FatherEnglishName,
					FatherAge:         f.FatherAge,
					FatherIsAlive:     f.FatherIsAlive,
					FatherPhoneNumber: f.FatherPhoneNumber,
					FatherOccupation:  f.FatherOccupation,
					FatherWorkplace:   f.FatherWorkplace,
					MotherName:        f.MotherName,
					MotherEnglishName: f.MotherEnglishName,
					MotherAge:         f.MotherAge,
					MotherIsAlive:     f.MotherIsAlive,
					MotherPhoneNumber: f.MotherPhoneNumber,
					MotherOccupation:  f.MotherOccupation,
					MotherWorkplace:   f.MotherWorkplace,
				})
			}
			if err := tx.Save(&family).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to update student family", nil)
			}
		}

		// Replace educations
		if err := tx.Where("student_id = ?", student.ID).
			Delete(&model.StudentEducation{}).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to clear student educations", nil)
		}
		if len(input.StudentEducationRequestUpdate) > 0 {
			educations := make([]model.StudentEducation, 0, len(input.StudentEducationRequestUpdate))
			for _, e := range input.StudentEducationRequestUpdate {
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

		// Replace documents
		if err := tx.Where("student_id = ?", student.ID).
			Delete(&model.StudentDocument{}).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to clear student documents", nil)
		}
		if len(input.StudentDocumentRequestUpdate) > 0 {
			documents := make([]model.StudentDocument, 0, len(input.StudentDocumentRequestUpdate))
			for _, d := range input.StudentDocumentRequestUpdate {
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

func (s *studentService) CreateStudent(ctx context.Context, input request.StudentRequestCreate) error {
	input.NameKh = strings.TrimSpace(input.NameKh)
	input.NameEn = strings.TrimSpace(input.NameEn)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Nationality = strings.TrimSpace(input.Nationality)
	username := strings.ToLower(input.NameEn)
	email := helper.GenerateEmail(username, 168)
	if err := studentValidator.Struct(input); err != nil {
		log.Printf("validation error: %v", err) // <-- add this
		return apperror.New(apperror.CodeInvalidInput, err.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	student := model.Student{
		GroupID:           input.GroupID,
		StudentCategoryID: input.StudentCategoryID,
		UserName:          username,
		Email:             email,
		Password:          utils.HasPassword("KUBB"),
		NameKh:            input.NameKh,
		NameEn:            input.NameEn,
		DateOfBirth:       input.DateOfBirth,
		Gender:            input.Gender,
		Nationality:       input.Nationality,
		Phone:             input.Phone,
		Status:            share.Created,
		VillageID:         &input.VillageID,
		Occupation:        input.Occupation,
		AcademicStreamID:  input.AcademicStreamID,
		TelegramUsername:  nil,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&student).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create student", nil)
		}

		student.Code = helper.GenerateCode("STU", uint(student.ID))

		if err := tx.Save(&student).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "faild to update code", nil)
		}

		if len(input.StudentFamilyRequestCreate) > 0 {
			family := make([]model.StudentFamily, 0, len(input.StudentFamilyRequestCreate))
			for _, f := range input.StudentFamilyRequestCreate {
				family = append(family, model.StudentFamily{
					StudentID:         student.ID,
					FatherName:        f.FatherName,
					FatherEnglishName: f.FatherEnglishName,
					FatherAge:         f.FatherAge,
					FatherIsAlive:     f.FatherIsAlive,
					FatherPhoneNumber: f.FatherPhoneNumber,
					FatherOccupation:  f.FatherOccupation,
					FatherWorkplace:   f.FatherWorkplace,
					MotherName:        f.MotherName,
					MotherEnglishName: f.MotherEnglishName,
					MotherAge:         f.MotherAge,
					MotherIsAlive:     f.MotherIsAlive,
					MotherPhoneNumber: f.MotherPhoneNumber,
					MotherOccupation:  f.MotherOccupation,
					MotherWorkplace:   f.MotherWorkplace,
				})
			}
			if err := tx.Create(&family).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create student family", nil)
			}
		}

		if input.AdmissionRequestCreate != nil {
			admission := model.Admission{
				UUIDBase: base.UUIDBase{
					UUID: helper.GenerateUUID(),
				},
				StudentID:        student.ID,
				TermID:           input.AdmissionRequestCreate.TermID,
				AcademicDegreeID: input.AdmissionRequestCreate.AcademicDegreeID,
				Date:             input.AdmissionRequestCreate.Date,
				State:            input.AdmissionRequestCreate.AdmissionState,
				Description:      input.AdmissionRequestCreate.Description,
				ReferralSchool:   input.AdmissionRequestCreate.ReferralSchool,
				Active:           true,
			}

			if err := tx.Create(&admission).Error; err != nil {
				return err
			}

			if input.EnrollmentRequestCreate != nil {
				enrollment := model.Enrollment{
					UUIDBase: base.UUIDBase{
						UUID: helper.GenerateUUID(),
					},
					AdmissionID:    admission.ID,
					SchoolarshipID: input.EnrollmentRequestCreate.SchoolarshipID,
					SectionID:      nil,
					FeeInterval:    input.EnrollmentRequestCreate.FeeInterval,
					Description:    nil,
				}
				if err := tx.Create(&enrollment).Error; err != nil {
					return err
				}

				var degree model.AcademicDegree
				if err := tx.First(&degree, input.AdmissionRequestCreate.AcademicDegreeID).Error; err != nil {
					return apperror.New(apperror.CodeInternal, "failed to load academic degree", nil)
				}

				baseAmount := helper.GetFeeAmountPerYear(degree, enrollment.FeeInterval)
				var discountgroupt *model.FeeDiscountGroup
				if student.GroupID > 0 {
					var group model.FeeDiscountGroup
					if err := tx.First(&group, student.GroupID).Error; err != nil {
						return apperror.New(apperror.CodeInternal, "failed to load discount group", nil)
					}
					discountgroupt = &group
				}

				discount := helper.CalculateDiscount(baseAmount, discountgroupt)
				total := baseAmount - discount
				if total < 0 {
					total = 0
				}
				var schoolarship model.Schoolarship
				if err := tx.First(&schoolarship, enrollment.SchoolarshipID).Error; err != nil {
					return apperror.New(apperror.CodeInternal, "failed to load discount group", nil)
				}
				secondDiscount := helper.CalculateDiscountBySchoolarship(baseAmount, &schoolarship)
				nettotal := total - secondDiscount
				totaldiscount := discount + secondDiscount
				fee := model.Fee{
					UUIDBase:     base.UUIDBase{UUID: helper.GenerateUUID()},
					EnrollmentID: enrollment.ID,
					Date:         time.Now().Format("2006-01-02"),
					Amount:       baseAmount,
					Discount:     totaldiscount,
					Total:        nettotal,
					Active:       true,
				}
				if err := tx.Create(&fee).Error; err != nil {
					return apperror.New(apperror.CodeInternal, "failed to create fee", nil)
				}

				scheduleCount, monthInterval := helper.GetNextDueDate(enrollment.FeeInterval)
				if secondDiscount > 0 {
					paymentAmount := nettotal / float64(scheduleCount)
					paymentAmount = math.Round(paymentAmount*100) / 100

					installments := make([]model.Installment, 0, scheduleCount)
					dueDate := time.Now()

					for i := 1; i <= scheduleCount; i++ {
						amount := paymentAmount
						if i == scheduleCount {
							paidSoFar := paymentAmount * float64(scheduleCount-1)
							amount = nettotal - paidSoFar
						}

						installments = append(installments, model.Installment{
							UUIDBase:   base.UUIDBase{UUID: helper.GenerateUUID()},
							FeeID:      fee.ID,
							SequenceNO: i,
							DueDate:    dueDate.Format("2006-01-02"),
							Amount:     amount,
							Status:     model.InstallmentStatusPending,
							InvoiceID:  nil,
						})
						dueDate = dueDate.AddDate(0, monthInterval, 0)
					}

					if err := tx.Create(&installments).Error; err != nil {
						return apperror.New(apperror.CodeInternal, "failed to create installments", nil)
					}
				}

				if input.StudentTermRequestCreate != nil {
					newstudentterm := model.StudentTerm{
						UUIDBase: base.UUIDBase{
							UUID: helper.GenerateUUID(),
						},
						EnrollmentID: enrollment.ID,
						SemesterID:   input.StudentTermRequestCreate.SemesterID,
						StudyYearID:  input.StudentTermRequestCreate.StudyYearID,
						Active:       true,
						Status:       "PENDING",
					}
					if err := tx.Create(&newstudentterm).Error; err != nil {
						return err
					}

					newgparecord := model.GpaRecord{
						UUIDBase: base.UUIDBase{
							UUID: helper.GenerateUUID(),
						},
						StudentTermID: newstudentterm.ID,
						TotalCredit:   0.00,
						SemesterGpa:   0.00,
						CumulativeGpa: 0.00,
					}
					if err := tx.Create(&newgparecord).Error; err != nil {
						return err
					}
				}

			}
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
			Joins("LEFT JOIN student_category sg ON sg.id = s.student_category_id").
			Joins("LEFT JOIN fee_discount_groups f ON f.id = s.group_id").
			Joins("LEFT JOIN academic_streams asd ON asd.id = s.academic_stream_id").
			Joins("LEFT JOIN villages sv ON sv.id = s.village_id").
			Joins("LEFT JOIN communes sc ON sc.id = sv.commune_id").
			Joins("LEFT JOIN districts sd ON sd.id = sc.district_id").
			Joins("LEFT JOIN provinces sp ON sp.id = sd.province_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where(
				"s.name_kh LIKE ? OR s.name_en LIKE ?",
				"%"+v+"%",
				"%"+v+"%",
			)
		}
		if v, ok := filter["student_category_id"]; ok && v != "" {
			tx = tx.Where("s.student_category_id = ?", v)
		}
		if v, ok := filter["group_id"]; ok && v != "" {
			tx = tx.Where("s.group_id = ?", v)
		}
		if v, ok := filter["phone"]; ok && v != "" {
			tx = tx.Where("s.phone LIKE ?", "%"+v+"%")
		}
		if v, ok := filter["stream_id"]; ok && v != "" {
			tx = tx.Where("s.academic_stream_id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).
		Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count students: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).
		Select(`
			sg.id AS student_category_id,
			sg.name AS student_category_name,
			s.id AS id,
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
			s.exam_in AS exam_in,
			s.exam_out AS exam_out,
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

	studentIDs := make([]int, 0, len(data))
	for _, student := range data {
		studentIDs = append(studentIDs, student.ID)
	}

	var studentFamilies []model.StudentFamily
	if err := s.db.WithContext(ctx).
		Table("student_families sf").
		Where("sf.student_id IN ?", studentIDs).
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
			sf.mother_occupation AS mother_occupation,
			sf.mother_workplace AS mother_workplace
		`).Scan(&studentFamilies).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch student family: %w", err)
	}

	var studentEducations []response.StudentEducationResponse
	if err := s.db.WithContext(ctx).
		Table("student_education sdu").
		Joins("LEFT JOIN villages sv ON sv.id = sdu.village_id").
		Joins("LEFT JOIN communes sc ON sc.id = sv.commune_id").
		Joins("LEFT JOIN districts sd ON sd.id = sc.district_id").
		Joins("LEFT JOIN provinces sp ON sp.id = sd.province_id").
		Where("sdu.student_id IN ?", studentIDs).
		Select(`
			sdu.id AS id,
			sdu.uuid AS uuid,
			sdu.student_id AS student_id,
			sdu.level AS level,
			sdu.school_name AS school_name,
			sv.id AS village_id,
			sv.name_kh AS villlage_name_kh,
			sc.id AS communce_id,
			sc.name_kh AS communce_name,
			sd.id AS district_id,
			sd.name_kh AS distirct_name,
			sp.id AS province_id,
			sp.name_kh AS province_name,
			sdu.start_date AS start_date,
			sdu.end_date AS end_date,
			sdu.cerificate_date AS cerificate_date,
			sdu.score AS score,
			sdu.gpa AS gpa,
			sdu.grade AS grade
		`).Scan(&studentEducations).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch student education: %w", err)
	}

	var studentDocuments []response.StudentDocumentResponse
	if err := s.db.WithContext(ctx).
		Table("student_documents sdo").
		Joins("LEFT JOIN document_types dt ON dt.id = sdo.document_type_id").
		Where("sdo.student_id IN ?", studentIDs).
		Select(`
			sdo.id AS id,
			sdo.uuid AS uuid,
			sdo.student_id AS student_id,
			dt.id AS document_type_id,
			dt.name_kh AS document_type_name_kh,
			sdo.required_qty AS required_qty,
			sdo.received_qty AS received_qty,
			sdo.remark AS remark
		`).Scan(&studentDocuments).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch student document: %w", err)
	}

	// group children by student_id
	familyByStudent := make(map[int][]model.StudentFamily, len(studentFamilies))
	for _, fam := range studentFamilies {
		familyByStudent[fam.StudentID] = append(familyByStudent[fam.StudentID], fam)
	}

	educationByStudent := make(map[int][]response.StudentEducationResponse, len(studentEducations))
	for _, edu := range studentEducations {
		educationByStudent[edu.StudentID] = append(educationByStudent[edu.StudentID], edu)
	}

	documentByStudent := make(map[int][]response.StudentDocumentResponse, len(studentDocuments))
	for _, doc := range studentDocuments {
		documentByStudent[doc.StudentID] = append(documentByStudent[doc.StudentID], doc)
	}

	for i := range data {
		id := data[i].ID
		data[i].StudentFamily = familyByStudent[id]
		data[i].StudentEducation = educationByStudent[id]
		data[i].StudentDocument = documentByStudent[id]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *studentService) GetCourseRegistration(ctx context.Context, filter map[string]string) ([]response.CourseRegistrationResponse, error) {
	var data []response.CourseRegistrationResponse

	query := s.db.WithContext(ctx).
		Table("course_registrations cr").
		Joins("LEFT JOIN student_terms st ON st.id = cr.student_term_id").
		Joins("LEFT JOIN enrollments e ON e.id = st.enrollment_id").
		Joins("LEFT JOIN admissions ad ON ad.id = e.admission_id").
		Joins("LEFT JOIN students stu ON stu.id = ad.student_id")

	if v, ok := filter["class_offering_id"]; ok && v != "" {
		query = query.Where("cr.class_offering_id = ?", v)
	}

	if v, ok := filter["attendance_id"]; ok && v != "" {
		query = query.Where(`
			NOT EXISTS (
				SELECT 1 FROM attendance_details adx
				WHERE adx.course_registration_id = cr.id
				AND adx.attendance_id = ?
			)
		`, v)
	}

	err := query.Select(`
			cr.id AS id,
			cr.uuid AS uuid,
			stu.name_kh AS name_kh,
			stu.name_en AS name_en,
			stu.date_of_birth AS date_of_birth,
			stu.gender AS gender,
			stu.phone AS phone
		`).
		Order("stu.id DESC").
		Scan(&data).Error
	if err != nil {
		return nil, fmt.Errorf("query course registrations: %w", err)
	}
	for i := range data {
		data[i].Dob = helper.FormatDate(data[i].Dob)
	}

	return data, nil
}
