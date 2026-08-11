package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/model/base"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TermService interface {
	GetTerm(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TermResponse, *model.PaginationMetadata, error)
	CreateTerm(ctx context.Context, input request.TermRequestCreate) error
	UpdateTerm(ctx context.Context, id string, input request.TermRequestUpdate) error
	Toggle(ctx context.Context, id string) error
}

type termservice struct {
	db *gorm.DB
}

func NewTermService() TermService {
	return &termservice{
		db: config.DB,
	}
}

func (s *termservice) GetTerm(
	ctx context.Context,
	pf request.Pagination,
	filter map[string]string,
) ([]response.TermResponse, *model.PaginationMetadata, error) {

	helper.NormalizePagination(&pf)

	var data []response.TermResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("terms t").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {

		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("g.id = ?", v)
		}

		if v, ok := filter["academic_id"]; ok && v != "" {
			tx = tx.Where("a.id = ?", v)
		}

		if v, ok := filter["active"]; ok && v != "" {
			tx = tx.Where("t.active = ?", v)
		}

		if v, ok := filter["search"]; ok && v != "" {
			search := "%" + v + "%"

			tx = tx.Where(`
				t.code LIKE ?
				OR t.name LIKE ?
			`, search, search)
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
			t.id AS id,
			t.uuid AS uuid,

			g.id AS generation_id,
			g.code AS generation_code,
			g.name AS generation_name,

			a.id AS academic_id,
			a.code AS academic_code,
			a.name AS academic_name,

			t.code AS code,
			t.name AS name,
			t.start_date AS start_date,
			t.end_date AS end_date,
			t.description AS description,
			t.active AS active
		`).
		Order("t.index ASC").
		Offset(offset).
		Limit(pf.PageSize)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	termIDs := make([]int, 0, len(data))

	for _, term := range data {
		termIDs = append(termIDs, term.ID)
	}

	type majorRow struct {
		MajorTermUUID   string `gorm:"column:major_term_uuid"`
		MajorTermActive bool   `gorm:"column:major_term_active"`
		TermID          int    `gorm:"column:term_id"`

		MajorID          int                    `gorm:"column:major_id"`
		MajorUUID        string                 `gorm:"column:major_uuid"`
		MajorCode        string                 `gorm:"column:major_code"`
		MajorName        string                 `gorm:"column:major_name"`
		DurationPeriod   int                    `gorm:"column:duration_period"`
		DurationInterval model.DurationInterval `gorm:"column:duration_interval"`
		MajorDescription string                 `gorm:"column:major_description"`
		MajorActive      bool                   `gorm:"column:major_active"`

		DepartmentID   int    `gorm:"column:department_id"`
		DepartmentName string `gorm:"column:department_name"`
		DepartmentCode string `gorm:"column:department_code"`

		FacultyID   int    `gorm:"column:faculty_id"`
		FacultyName string `gorm:"column:faculty_name"`
		FacultyCode string `gorm:"column:faculty_code"`

		ProgrammeID   int    `gorm:"column:programme_id"`
		ProgrammeName string `gorm:"column:programme_name"`
	}

	var majorRows []majorRow

	if err := s.db.WithContext(ctx).
		Table("major_terms mt").
		Joins("INNER JOIN majors m ON m.id = mt.major_id").
		Joins("LEFT JOIN departments d ON d.id = m.department_id").
		Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
		Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
		Where("mt.term_id IN ?", termIDs).
		Select(`
			mt.term_id AS term_id,
			mt.uuid AS major_term_uuid,
			mt.active AS major_term_active,
			m.id AS major_id,
			m.uuid AS major_uuid,
			m.code AS major_code,
			m.name AS major_name,
			m.duration_period AS duration_period,
			m.duration_interval AS duration_interval,
			m.description AS major_description,
			m.active AS major_active,

			d.id AS department_id,
			d.name AS department_name,
			d.code AS department_code,

			f.id AS faculty_id,
			f.name AS faculty_name,
			f.code AS faculty_code,

			p.id AS programme_id,
			p.name AS programme_name
		`).
		Scan(&majorRows).Error; err != nil {

		return nil, nil, fmt.Errorf("fetch term majors: %w", err)
	}

	majorsByTerm := make(map[int][]response.MajorResponse)

	for _, row := range majorRows {

		major := response.MajorResponse{
			MajorTermUUID:    row.MajorTermUUID,
			MajorTermActive:  row.MajorTermActive,
			ID:               row.MajorID,
			UUID:             row.MajorUUID,
			Name:             row.MajorName,
			Code:             row.MajorCode,
			DurationPeriod:   row.DurationPeriod,
			DurationInterval: row.DurationInterval,
			Description:      row.MajorDescription,
			Active:           row.MajorActive,

			DepartmentID:   row.DepartmentID,
			DepartmentName: row.DepartmentName,
			DepartmentCode: row.DepartmentCode,

			FacultyID:   row.FacultyID,
			FacultyName: row.FacultyName,
			FacultyCode: row.FacultyCode,

			ProgrammeID:  row.ProgrammeID,
			ProgrammName: row.ProgrammeName,
		}

		majorsByTerm[row.TermID] = append(
			majorsByTerm[row.TermID],
			major,
		)
	}

	for i := range data {
		data[i].MajorResponse = majorsByTerm[data[i].ID]
		if data[i].MajorResponse == nil {
			data[i].MajorResponse = []response.MajorResponse{}
		}
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *termservice) CreateTerm(ctx context.Context, input request.TermRequestCreate) error {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation code is required", nil)
	}

	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "generation name is required", nil)
	}

	if input.GenerationID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "generation id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	nextIndex := 1

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lastTerm model.Term
		err := tx.
			Where("generation_id = ?", input.GenerationID).
			Order("id DESC").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lastTerm).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			nextIndex = 1
		case err != nil:
			return helper.MapAcademicError(err, "LOOKUP_LAST_GENERATION")
		default:
			nextIndex = lastTerm.Index + 1
		}

		newdata := model.Term{
			UUIDBase: base.UUIDBase{
				UUID: helper.GenerateUUID(),
			},
			GenerationID: input.GenerationID,
			Code:         code,
			Name:         name,
			Index:        nextIndex,
			StartDate:    input.StartDate,
			EndDate:      input.EndDate,
			Description:  input.Description,
			Active:       true,
		}

		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapAcademicError(err, "CREATE")
		}
		return nil
	})
	return err
}

func (s *termservice) UpdateTerm(ctx context.Context, id string, input request.TermRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "id is required", nil)
	}

	updates := map[string]interface{}{}

	if input.GenerationID != nil {
		updates["generation_id"] = *input.GenerationID
	}

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return apperror.New(apperror.CodeInvalidInput, "generation code is required", nil)
		}
		updates["code"] = code
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "generation name is required", nil)
		}
		updates["name"] = name
	}

	if input.StartDate != nil {
		updates["start_date"] = *input.StartDate
	}

	if input.EndDate != nil {
		updates["end_date"] = *input.EndDate
	}

	if input.Description != nil {
		updates["description"] = *input.Description
	}

	if len(updates) == 0 {
		return apperror.Invalid("no fields provided to update", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).
		Model(&model.Term{}).
		Where("uuid = ?", id).
		Updates(updates)

	if result.Error != nil {
		return helper.MapAcademicError(result.Error, "update")
	}

	if result.RowsAffected == 0 {
		return apperror.NotFound("generation not found", nil)
	}

	return nil
}

func (s *termservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Term](ctx, s.db, id)
}
