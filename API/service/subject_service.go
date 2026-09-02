package service

import (
	"context"
	"errors"
	"fmt"
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

type SubjectService interface {
	CreateSubject(ctx context.Context, input request.SubjectRequestCreate) error
	GetSubject(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.SubjectResponse, *model.PaginationMetadata, error)
	GetSubjectByMajor(ctx context.Context, majorID int) ([]response.SubjectResponseByMajor, error)
	Toggle(ctx context.Context, id string) error
	UpdateSubject(ctx context.Context, id string, input request.SubjectRequestUpdate) error
	CreateGradeComponent(ctx context.Context, input request.GradeComponentRequestCreate) error
	GetSubjectGroup(ctx context.Context) ([]model.SubjectGroup, error)
}

type subjectservice struct {
	db *gorm.DB
}

func NewSubjectService() SubjectService {
	return &subjectservice{
		db: config.DB,
	}
}

func (s *subjectservice) GetSubjectGroup(ctx context.Context) ([]model.SubjectGroup, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var data []model.SubjectGroup
	err := s.db.WithContext(ctx).Find(&data).Error
	if err != nil {
		return nil, apperror.Internal("failed to fetch subject group", err)
	}
	return data, nil
}

func (s *subjectservice) CreateSubject(ctx context.Context, input request.SubjectRequestCreate) error {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)

	if code == "" {
		return apperror.New(apperror.CodeInvalidInput, "subject code is required", nil)
	}
	if name == "" {
		return apperror.New(apperror.CodeInvalidInput, "subject name is required", nil)
	}
	if input.MajorID == 0 {
		return apperror.New(apperror.CodeInvalidInput, "major_id is required", nil)
	}
	if input.Credit < 0 {
		return apperror.New(apperror.CodeInvalidInput, "credit must be zero or positive", nil)
	}
	if input.PassingScore < 0 {
		return apperror.New(apperror.CodeInvalidInput, "passing_score must be zero or positive", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	code = strings.ToUpper(code)

	newdata := model.Subject{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		MajorID:      input.MajorID,
		Code:         code,
		Name:         name,
		Credit:       input.Credit,
		PassingScore: input.PassingScore,
		Description:  description,
		Active:       true,
	}

	var count int64
	if err := s.db.WithContext(ctx).
		Model(&model.Subject{}).
		Where("major_id = ? AND code = ?", input.MajorID, code).
		Count(&count).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	if count > 0 {
		return apperror.New(apperror.CodeConflict, "subject code already exists for this major", nil)
	}

	if err := s.db.WithContext(ctx).Create(&newdata).Error; err != nil {
		return helper.MapAcademicError(err, "CREATE")
	}

	return nil
}

func (s *subjectservice) GetSubject(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.SubjectResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.SubjectResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("subjects s").
			Joins("LEFT JOIN majors m ON m.id = s.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["programme_id"]; ok && v != "" {
			tx = tx.Where("p.id = ?", v)
		}
		if v, ok := filter["faculty_id"]; ok && v != "" {
			tx = tx.Where("f.id = ?", v)
		}
		if v, ok := filter["department_id"]; ok && v != "" {
			tx = tx.Where("d.id = ?", v)
		}
		if v, ok := filter["major_id"]; ok && v != "" {
			tx = tx.Where("m.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count terms: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).Select(`
		s.id AS id,
		s.uuid AS uuid,
		s.code AS code,
		s.name AS name,
		s.credit AS credit,
		s.passing_score AS passing_score,
		s.description AS description,
		s.active AS active,
		m.id AS major_id,
		m.code AS major_code,
		m.name AS major_name,
		d.id AS department_id,
		d.name AS department_name,
		d.code AS department_code,
		f.id AS faculty_id,
		f.code AS faculty_code,
		f.name AS faculty_name,
		p.id AS programme_id,
		p.name AS programme_name
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch terms: %w", err)
	}

	subjectIDs := make([]int, 0, len(data))
	for _, s := range data {
		subjectIDs = append(subjectIDs, s.ID)
	}

	var gradeComponents []response.GradeComponentResponse
	if len(subjectIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("grade_components gc").
			Where("gc.subject_id IN ?", subjectIDs).
			Select(`
				gc.id AS id,
				gc.uuid AS uuid,
				gc.subject_id AS subject_id,
				gc.name AS name,
				gc.weight_percentage AS weight_percentage,
				gc.active AS active
			`).
			Scan(&gradeComponents).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch grade components: %w", err)
		}
	}
	gradeComponentsBySubject := make(map[int][]response.GradeComponentResponse)
	for _, gc := range gradeComponents {
		gradeComponentsBySubject[gc.SubjectID] =
			append(gradeComponentsBySubject[gc.SubjectID], gc)
	}
	for i := range data {
		data[i].GradeComponentResponse =
			gradeComponentsBySubject[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *subjectservice) GetSubjectByMajor(ctx context.Context, majorID int) ([]response.SubjectResponseByMajor, error) {
	if majorID <= 0 {
		return nil, apperror.New(apperror.CodeInvalidInput, "major id  is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.SubjectResponseByMajor

	err := s.db.WithContext(ctx).
		Table("subjects s").
		Select(`
			s.id AS id,
			s.name AS name
		`).
		Where("s.major_id = ?", majorID).
		Find(&data).Error

	if err != nil {
		return nil, fmt.Errorf("fetch major by major: %w", err)
	}

	return data, nil
}

func (s *subjectservice) Toggle(ctx context.Context, id string) error {
	return utils.ToggleStatus[model.Subject](ctx, s.db, id)
}

func (s *subjectservice) UpdateSubject(ctx context.Context, id string, input request.SubjectRequestUpdate) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "subject id is required", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var existing model.Subject
	if err := s.db.WithContext(ctx).
		Where("uuid = ?", id).
		First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.CodeNotFound, "subject not found", nil)
		}
		return helper.MapAcademicError(err, "UPDATE")
	}

	updates := map[string]interface{}{}

	if input.MajorID != nil {
		if *input.MajorID == 0 {
			return apperror.New(apperror.CodeInvalidInput, "major_id cannot be zero", nil)
		}
		updates["major_id"] = *input.MajorID
	}

	if input.Code != nil {
		code := strings.ToUpper(strings.TrimSpace(*input.Code))
		if code == "" {
			return apperror.New(apperror.CodeInvalidInput, "subject code cannot be empty", nil)
		}
		updates["code"] = code
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "subject name cannot be empty", nil)
		}
		updates["name"] = name
	}

	if input.Credit != nil {
		if *input.Credit < 0 {
			return apperror.New(apperror.CodeInvalidInput, "credit must be zero or positive", nil)
		}
		updates["credit"] = *input.Credit
	}

	if input.PassingScore != nil {
		if *input.PassingScore < 0 {
			return apperror.New(apperror.CodeInvalidInput, "passing_score must be zero or positive", nil)
		}
		updates["passing_score"] = *input.PassingScore
	}

	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}

	if len(updates) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "no fields provided to update", nil)
	}

	effectiveMajorID := existing.MajorID
	if v, ok := updates["major_id"].(int); ok {
		effectiveMajorID = v
	}
	effectiveCode := existing.Code
	if v, ok := updates["code"].(string); ok {
		effectiveCode = v
	}

	if effectiveMajorID != existing.MajorID || effectiveCode != existing.Code {
		var count int64
		if err := s.db.WithContext(ctx).
			Model(&model.Subject{}).
			Where("major_id = ? AND code = ? AND uuid <> ?", effectiveMajorID, effectiveCode, existing.UUID).
			Count(&count).Error; err != nil {
			return helper.MapAcademicError(err, "UPDATE")
		}
		if count > 0 {
			return apperror.New(apperror.CodeConflict, "subject code already exists for this major", nil)
		}
	}

	if err := s.db.WithContext(ctx).
		Model(&existing).
		Updates(updates).Error; err != nil {
		return helper.MapAcademicError(err, "UPDATE")
	}

	return nil
}

func (s *subjectservice) CreateGradeComponent(ctx context.Context, input request.GradeComponentRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(input.GradeComponentRequest) > 0 {
			data := make([]model.GradeComponent, 0, len(input.GradeComponentRequest))
			for _, g := range input.GradeComponentRequest {
				data = append(data, model.GradeComponent{
					UUIDBase: base.UUIDBase{
						UUID: helper.GenerateUUID(),
					},
					SubjectID:        input.SubjectID,
					Name:             g.Name,
					WeightPercentage: g.WeightPercentage,
					Active:           true,
				})
			}
			if err := tx.Create(&data).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create grade component", nil)
			}
		}
		return nil
	})
	return err
}
