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
			Joins("fee_discount_groups fd ON fd.id = s.group_id").
			Joins("LEFT JOIN terms t ON t.id = adm.term_id").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id").
			Joins("LEFT JOIN academic_degrees ad ON ad.id = adm.academic_degree_id").
			Joins("LEFT JOIN majors m ON m.id = adm.major_id").
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

	if err := applyFilters(base()).
		Count(&total).Error; err != nil {

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
			adm.student_id AS student_id,
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
			ad.yearly_fee AS yearly_fee,		
		`).
		Offset(offset).
		Limit(pf.PageSize)
	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}
	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

}
