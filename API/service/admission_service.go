package service

import (
	"context"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type AdmissionService interface {
	GetAdmission(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.AdmissionResponse, *model.PaginationMetadata, error)
}

type admissionservice struct {
	db *gorm.DB
}

func NewAdmissionService() AdmissionService {
	return &admissionservice{
		db: config.DB,
	}
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
			tx = tx.Where("s.id = ?", v)
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
			s.name_kh AS student_name_kh,
			s.name_en AS student_name_en,
			s.gender AS student_gender,
			fd.discount_type AS discount_type,
			fd.discount_percentage AS discount_percentage,
			fd.discount_amount AS discount_amount,
			t.id AS term_id,
			t.name AS term_name,
			g.code AS generation_code,
			g.name AS generation_name,
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
		Where("e.admission_id IN ?", admissionIDs).
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
			e.description AS description
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

	var invoices []response.InvoiceResposne
	if len(feeIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("invoices i").
			Joins("INNER JOIN fees f ON f.id = i.fee_id").
			Where("i.fee_id IN ?", feeIDs).
			Select(`
				i.id AS id,
				i.uuid AS uuid,
				i.fee_id AS fee_id,
				i.code AS code,
				i.invoice_date AS invoice_date,
				i.due_date AS due_date,
				i.total AS total,
				i.discount AS discount,
				i.tax AS tax,
				i.grant_total AS grant_total,
				i.message_on_invoice AS message_on_invoice,
				i.description AS description,
				i.active AS active
			`).Scan(&invoices).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch invoices: %w", err)
		}
	}

	var installments []response.InstallmentResponse
	if len(feeIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("installments it").
			Joins("INNER JOIN fees f ON f.id = it.fee_id").
			Where("it.fee_id IN ?", feeIDs).
			Select(`
				it.id AS id,
				it.uuid AS uuid,
				it.fee_id AS fee_id,
				it.sequence_no AS sequence_no,
				it.due_date AS due_date,
				it.amount AS amount,
				it.status AS status,
				it.invoice_id AS invoice_id
			`).Scan(&installments).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch installments: %w", err)
		}
	}

	for i := range installments {
		installments[i].DueDate = helper.FormatDate(installments[i].DueDate)
	}

	invoiceIDs := make([]int, 0, len(invoices))
	for _, invoice := range invoices {
		invoiceIDs = append(invoiceIDs, invoice.ID)
	}

	var payments []response.PaymentResposen
	if len(invoiceIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("payments p").
			Joins("INNER JOIN invoices i ON i.id = p.invoice_id").
			Where("p.invoice_id IN ?", invoiceIDs).
			Select(`
				p.id AS id,
				p.uuid AS uuid,
				p.invoice_id AS invoice_id,
				p.code AS code,
				p.date AS date,
				p.amount AS amount,
				p.reference AS reference,
				p.method AS method,
				p.description AS description,
				p.active AS active
			`).Scan(&payments).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch payments: %w", err)
		}
	}

	// ---- Assemble the nested tree ----

	paymentsByInvoice := make(map[int][]response.PaymentResposen, len(payments))
	for _, pay := range payments {
		paymentsByInvoice[pay.InvoiceID] = append(paymentsByInvoice[pay.InvoiceID], pay)
	}

	invoicesByFee := make(map[int][]response.InvoiceResposne, len(invoices))
	for _, inv := range invoices {
		inv.PaymentResposen = paymentsByInvoice[inv.ID]
		invoicesByFee[inv.FeeID] = append(invoicesByFee[inv.FeeID], inv)
	}

	installmentsByFee := make(map[int][]response.InstallmentResponse, len(installments))
	for _, inst := range installments {
		if inst.InvoiceID != nil {
			if invs := invoicesByFee[inst.FeeID]; len(invs) > 0 {
				for _, inv := range invs {
					if inv.ID == *inst.InvoiceID {
						inst.InvoiceResposne = inv
						break
					}
				}
			}
		}
		installmentsByFee[inst.FeeID] = append(installmentsByFee[inst.FeeID], inst)
	}

	feesByEnrollment := make(map[int][]response.FeeResponse, len(fees))
	for _, fee := range fees {
		fee.InvoiceResposne = invoicesByFee[fee.ID]
		fee.InstallmentResponse = installmentsByFee[fee.ID]
		feesByEnrollment[fee.EnrollmentID] = append(feesByEnrollment[fee.EnrollmentID], fee)
	}

	studentTermsByEnrollment := make(map[int][]response.StudentTermResponse, len(studentTerms))
	for _, st := range studentTerms {
		st.FeeResponse = feesByEnrollment[st.EnrollmentID]
		studentTermsByEnrollment[st.EnrollmentID] = append(studentTermsByEnrollment[st.EnrollmentID], st)
	}

	enrollmentsByAdmission := make(map[int][]response.EnrollmentResponse, len(enrollments))
	for _, en := range enrollments {
		en.StudentResponse = studentTermsByEnrollment[en.ID]
		enrollmentsByAdmission[en.AdmissionID] = append(enrollmentsByAdmission[en.AdmissionID], en)
	}

	for i := range data {
		data[i].EnrollmentResponse = enrollmentsByAdmission[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
