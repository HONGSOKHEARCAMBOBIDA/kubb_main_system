package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"time"

	"gorm.io/gorm"
)

type AdmissionService interface {
	GetAdmission(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AdmissionResponse, *model.PaginationMetadata, error)
	CreateStudentTerm(ctx context.Context, input request.StudentTermRequestv2) error
	CreateEnrollment(ctx context.Context, input request.EnrollmentRequestCreateV2) error
	GetStudentTermFilter(ctx context.Context, filter map[string]string) ([]response.StudentTermResponsebyFilter, error)
	UpdateAdmission(ctx context.Context, uuid string, input request.AdmissionRequestUpdate) error
	UpdateEnrollment(ctx context.Context, uuid string, input request.EnrollmentRequestUpdate) error
	DeleteEnrollment(ctx context.Context, uuid string) error
	UpdateStudentTerm(ctx context.Context, uuid string, input request.StudentTermRequestUpdate) error
}

type admissionservice struct {
	db *gorm.DB
}

func NewAdmissionService() AdmissionService {
	return &admissionservice{
		db: config.DB,
	}
}

func (s *admissionservice) DeleteEnrollment(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var enrollment model.Enrollment
		if err := tx.Where("uuid = ?", uuid).First(&enrollment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "enrollment not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch enrollment", nil)
		}
		enrollment.Isactive = false
		if err := tx.Save(&enrollment).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *admissionservice) UpdateStudentTerm(ctx context.Context, uuid string, input request.StudentTermRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var studentterm model.StudentTerm
		if err := tx.Where("uuid = ?", uuid).First(&studentterm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "studentterm not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch studentterm", nil)
		}
		studentterm.SemesterID = input.SemesterID
		studentterm.StudyYearID = input.StudyYearID
		studentterm.Status = input.Status
		if err := tx.Save(&studentterm).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *admissionservice) UpdateEnrollment(ctx context.Context, uuid string, input request.EnrollmentRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var enrollment model.Enrollment
		if err := tx.Where("uuid = ?", uuid).First(&enrollment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "enrollment not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch enrollment", nil)
		}
		enrollment.Description = &input.Description
		if err := tx.Save(&enrollment).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *admissionservice) UpdateAdmission(ctx context.Context, uuid string, input request.AdmissionRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var admission model.Admission
		if err := tx.Where("uuid = ?", uuid).First(&admission).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "admission not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch admission", nil)
		}
		admission.TermID = input.TermID
		admission.State = input.AdmissionState
		if err := tx.Save(&admission).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *admissionservice) GetStudentTermFilter(ctx context.Context, filter map[string]string) ([]response.StudentTermResponsebyFilter, error) {
	var data []response.StudentTermResponsebyFilter

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("student_terms st").
			Joins("LEFT JOIN semesters ss ON ss.id = st.semester_id").
			Joins("LEFT JOIN enrollments e ON e.id = st.enrollment_id").
			Joins("LEFT JOIN admissions adm ON adm.id = e.admission_id").
			Joins("LEFT JOIN students s ON s.id = adm.student_id").
			Joins("LEFT JOIN terms t ON t.id = adm.term_id").
			Joins("LEFT JOIN academic_degrees ad ON ad.id = adm.academic_degree_id").
			Joins("LEFT JOIN majors m ON m.id = ad.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["semester_id"]; ok && v != "" {
			tx = tx.Where("st.semester_id = ?", v)
		}
		if v, ok := filter["study_year_id"]; ok && v != "" {
			tx = tx.Where("st.study_year_id = ?", v)
		}
		if v, ok := filter["term_id"]; ok && v != "" {
			tx = tx.Where("t.id = ?", v)
		}
		if v, ok := filter["major_id"]; ok && v != "" {
			tx = tx.Where("m.id = ?", v)
		}
		tx = tx.Where("st	.status = ?", "PENDING")
		// tx = tx.Where(`
		// 	NOT EXISTS (
		// 		SELECT 1
		// 		FROM course_registrations cr
		// 		WHERE cr.student_term_id = st.id
		// 	)
		// `)
		tx = tx.Where("e.is_active = ?", true)

		return tx
	}

	dataQuery := applyFilters(base()).
		Select(`
			st.id AS id,
			st.uuid AS uuid,
			s.id AS student_id,
			s.name_kh AS student_name_kh,
			s.name_en AS student_name_en,
			s.gender AS student_gender,
			ss.id AS semester_id,
			ss.name AS semester_name,
			st.study_year_id AS study_year_id,
			t.id AS term_id,
			t.code AS term_code,
			t.name AS term_name,
			m.id AS major_id,
			m.code AS major_code,
			m.name AS major_name,
			p.id AS program_id,
			p.name AS programm_name
		`)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, fmt.Errorf("fetch terms: %w", err)
	}

	return data, nil
}

func (s *admissionservice) GetAdmission(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AdmissionResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.AdmissionResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("admissions adm").
			Joins("LEFT JOIN students s ON s.id = adm.student_id").
			Joins("LEFT JOIN fee_discount_groups fd ON fd.id = s.group_id").
			Joins("LEFT JOIN terms t ON t.id = adm.term_id").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id").
			Joins("LEFT JOIN academic_degrees ad ON ad.id = adm.academic_degree_id").
			Joins("LEFT JOIN majors m ON m.id = ad.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["student_id"]; ok && v != "" {
			tx = tx.Where("s.code = ?", v)
		}
		if v, ok := filter["academic_id"]; ok && v != "" {
			tx = tx.Where("a.id = ?", v)
		}
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("g.id = ?", v)
		}
		if v, ok := filter["term_id"]; ok && v != "" {
			tx = tx.Where("t.id = ?", v)
		}
		if v, ok := filter["student_name"]; ok && v != "" {
			tx = tx.Where(
				"s.name_kh LIKE ? OR s.name_en LIKE ?",
				"%"+v+"%",
				"%"+v+"%",
			)
		}

		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count data: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).
		Select(`
			adm.id AS id,
			adm.uuid AS uuid,
			adm.date AS date,
			adm.state AS state,
			adm.description AS description,
			adm.referral_school AS referral_school,
			adm.active AS active,
			s.id AS student_id,
			s.code AS student_code,
			s.name_kh AS student_name_kh,
			s.name_en AS student_name_en,
			s.gender AS student_gender,
			fd.name AS group_name,
			fd.discount_type AS discount_type,
			fd.discount_percentage AS discount_percentage,
			fd.discount_amount AS discount_amount,
			t.id AS term_id,
			t.name AS term_name,
			g.id AS generation_id,
			g.code AS generation_code,
			g.name AS generation_name,
			a.id AS academic_id,
			a.code AS academic_code,
			a.name AS academic_name,
			ad.id AS academic_degree_id,
			m.code AS major_code,
			m.name AS major_name,
			p.id AS programme_id,
			p.name AS programme_name,
			ad.monthly_fee AS monthly_fee,
			ad.quarterly_fee AS quarterly_fee,
			ad.semesterly_fee AS semesterly_fee,
			ad.yearly_fee AS yearly_fee
		`).
		Offset(offset).
		Limit(pf.PageSize)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch admissions: %w", err)
	}
	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	admissionIDs := make([]int, 0, len(data))
	for _, admission := range data {
		admissionIDs = append(admissionIDs, admission.ID)
	}

	var enrollments []response.EnrollmentResponse
	if err := s.db.WithContext(ctx).
		Table("enrollments e").
		Joins("INNER JOIN admissions a ON a.id = e.admission_id").
		Joins("LEFT JOIN scholarships s ON s.id = e.scholarship_id").
		Joins("LEFT JOIN student_terms st ON st.id = e.id").
		Where("e.admission_id IN ? AND e.is_active = ?", admissionIDs, true).
		Select(`
			e.id AS id,
			e.uuid AS uuid,
			e.admission_id AS admission_id,
			s.id AS schoolarship_id,
			s.name AS schoolarship_name,
			s.discount_type AS schoolarship_discount_type,
			s.discount_amount AS schoolarship_discount_amount,
			s.discount_percentage AS schoolarship_discount_percentage,
			e.fee_interval AS fee_interval,
			e.description AS description,
			(
				SELECT st.study_year_id
				FROM student_terms st
				WHERE st.enrollment_id = e.id
				LIMIT 1
			) AS year_id
		`).Scan(&enrollments).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch enrollments: %w", err)
	}

	enrollmentIDs := make([]int, 0, len(enrollments))
	for _, enrollment := range enrollments {
		enrollmentIDs = append(enrollmentIDs, enrollment.ID)
	}

	var studentTerms []response.StudentTermResponse
	if len(enrollmentIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("student_terms st").
			Joins("INNER JOIN enrollments e ON e.id = st.enrollment_id").
			Joins("LEFT JOIN semesters s ON s.id = st.semester_id").
			Joins("LEFT JOIN academics a ON a.id = s.academic_id").
			Where("st.enrollment_id IN ?", enrollmentIDs).
			Select(`
				st.id AS id,
				st.uuid AS uuid,
				st.enrollment_id AS enrollment_id,
				s.id AS semester_id,
				s.code AS semester_code,
				s.name AS semester_name,
				a.id AS academic_id,
				a.name AS academic_name,
				st.study_year_id AS study_year_id,
				st.active AS active,
				st.status AS status
			`).Scan(&studentTerms).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch student terms: %w", err)
		}
	}

	studentTermIDs := make([]int, 0, len(studentTerms))
	for _, studentterm := range studentTerms {
		studentTermIDs = append(studentTermIDs, studentterm.ID)
	}

	var gparesponse []response.GpaRecordResponse
	if len(studentTermIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("gpa_records gr").
			Where("gr.student_term_id IN ?", studentTermIDs).
			Select(`
			gr.id AS id,
			gr.uuid AS uuid,
			gr.student_term_id AS student_term_id,
			gr.total_credit AS total_credit,
			gr.semester_gpa AS semester_gpa,
			gr.cumulative_gpa AS cumulative_gpa
		`).Scan(&gparesponse).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch student terms: %w", err)
		}
	}

	var fees []response.FeeResponse
	if len(enrollmentIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("fees f").
			Joins("INNER JOIN enrollments e ON e.id = f.enrollment_id").
			Where("f.enrollment_id IN ?", enrollmentIDs).
			Select(`
				f.id AS id,
				f.uuid AS uuid,
				f.enrollment_id AS enrollment_id,
				f.date AS date,
				f.amount AS amount,
				f.discount AS discount,
				f.total AS total,
				f.active AS active
			`).Scan(&fees).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch fees: %w", err)
		}
	}

	for i := range fees {
		fees[i].Date = helper.FormatDate(fees[i].Date)
	}

	feeIDs := make([]int, 0, len(fees))
	for _, fee := range fees {
		feeIDs = append(feeIDs, fee.ID)
	}

	var installments []response.InstallmentResponse
	if len(feeIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("installments it").
			Joins("INNER JOIN fees f ON f.id = it.fee_id").
			Joins("LEFT JOIN invoices i ON i.id = it.invoice_id").
			Joins("LEFT JOIN payments p ON p.invoice_id = i.id").
			Where("it.fee_id IN ?", feeIDs).
			Select(`
				it.id AS id,
				it.uuid AS uuid,
				it.fee_id AS fee_id,
				it.sequence_no AS sequence_no,
				it.due_date AS due_date,
				it.amount AS amount,
				it.status AS status,
				i.id AS invoice_id,
				i.code AS invoice_code,
				i.invoice_date AS invoice_date,
				i.due_date AS invoice_due_date,
				i.total AS invoice_total,
				i.tax AS invoice_tax,
				i.grant_total AS invoice_grant_total,
				i.message_on_invoice AS message_on_invoice,
				p.id AS payment_id,
				p.code AS payment_code,
				p.reference AS payment_reference,
				p.method AS payment_method
			`).Scan(&installments).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch installments: %w", err)
		}
	}

	for i := range installments {
		installments[i].DueDate = helper.FormatDate(installments[i].DueDate)
		installments[i].InvoiceDate = helper.FormatDate(installments[i].InvoiceDate)
	}

	installmentsByFee := make(map[int][]response.InstallmentResponse, len(installments))
	for _, inst := range installments {

		installmentsByFee[inst.FeeID] = append(installmentsByFee[inst.FeeID], inst)
	}

	feesByEnrollment := make(map[int][]response.FeeResponse, len(fees))
	for _, fee := range fees {
		fee.InstallmentResponse = installmentsByFee[fee.ID]
		feesByEnrollment[fee.EnrollmentID] = append(feesByEnrollment[fee.EnrollmentID], fee)
	}

	gpaByStudentTerm := make(map[int][]response.GpaRecordResponse, len(gparesponse))
	for _, gpa := range gparesponse {
		gpaByStudentTerm[gpa.StudentTermID] = append(gpaByStudentTerm[gpa.StudentTermID], gpa)
	}

	studentTermsByEnrollment := make(map[int][]response.StudentTermResponse, len(studentTerms))
	for _, st := range studentTerms {
		st.GpaRecordResponse = gpaByStudentTerm[st.ID]
		studentTermsByEnrollment[st.EnrollmentID] = append(studentTermsByEnrollment[st.EnrollmentID], st)
	}

	enrollmentsByAdmission := make(map[int][]response.EnrollmentResponse, len(enrollments))
	for _, en := range enrollments {
		en.FeeResponse = feesByEnrollment[en.ID]
		en.StudentResponse = studentTermsByEnrollment[en.ID]
		enrollmentsByAdmission[en.AdmissionID] = append(enrollmentsByAdmission[en.AdmissionID], en)
	}

	for i := range data {
		data[i].EnrollmentResponse = enrollmentsByAdmission[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *admissionservice) CreateStudentTerm(
	ctx context.Context,
	input request.StudentTermRequestv2,
) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		studentTerm := model.StudentTerm{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			EnrollmentID: input.EnrollmentID,
			SemesterID:   input.SemesterID,
			StudyYearID:  input.StudyYearID,
			Active:       true,
			Status:       "PENDING",
		}

		if err := tx.Create(&studentTerm).Error; err != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to create student term",
				nil,
			)
		}
		if err := tx.Model(&model.StudentTerm{}).
			Where("enrollment_id = ?", input.EnrollmentID).
			Where("id < ?", studentTerm.ID).
			Update("status", "FINISH").Error; err != nil {
			return apperror.New(
				apperror.CodeInternal,
				"failed to finish previous student terms",
				nil,
			)
		}
		newgparecord := model.GpaRecord{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			StudentTermID: studentTerm.ID,
			TotalCredit:   0.00,
			SemesterGpa:   0.00,
			CumulativeGpa: 0.00,
		}
		if err := tx.Create(&newgparecord).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *admissionservice) CreateEnrollment(ctx context.Context, input request.EnrollmentRequestCreateV2) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	enrollment := model.Enrollment{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		AdmissionID:    input.AdmissionID,
		SchoolarshipID: input.SchoolarshipID,
		FeeInterval:    input.FeeInterval,
		Isactive:       true,
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var admission model.Admission
		if err := tx.First(&admission, input.AdmissionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "admission not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to load admission", nil)
		}

		var student model.Student
		if err := tx.First(&student, admission.StudentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "student not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to load student", nil)
		}

		var degree model.AcademicDegree
		if err := tx.First(&degree, admission.AcademicDegreeID).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to load academic degree", nil)
		}

		scheduleCount := helper.GetFeeSchedule(enrollment.FeeInterval)
		if scheduleCount <= 0 {
			return apperror.New(apperror.CodeInternal, "invalid fee schedule for given interval", nil)
		}

		if err := tx.Create(&enrollment).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create enrollment", nil)
		}

		baseAmount := helper.GetFeeAmountPerYear(degree, enrollment.FeeInterval)

		var discountGroup *model.FeeDiscountGroup
		if student.GroupID > 0 {
			var group model.FeeDiscountGroup
			if err := tx.First(&group, student.GroupID).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to load discount group", nil)
			}
			discountGroup = &group
		}
		discount := helper.CalculateDiscount(baseAmount, discountGroup)
		total := baseAmount - discount
		if total < 0 {
			total = 0
		}

		var secondDiscount float64
		if enrollment.SchoolarshipID > 0 {
			var schoolarship model.Schoolarship
			if err := tx.First(&schoolarship, enrollment.SchoolarshipID).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to load schoolarship", nil)
			}
			secondDiscount = helper.CalculateDiscountBySchoolarship(baseAmount, &schoolarship)
		}

		nettotal := total - secondDiscount
		if nettotal < 0 {
			nettotal = 0
		}
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

		// --- Installments: always generated, regardless of discounts ---
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
			dueDate = dueDate.AddDate(0, 1, 0)
		}

		if err := tx.Create(&installments).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create installments", nil)
		}

		// --- Optional student term ---
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
				return apperror.New(apperror.CodeInternal, "failed to create student term", nil)
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

		return nil
	})
}
