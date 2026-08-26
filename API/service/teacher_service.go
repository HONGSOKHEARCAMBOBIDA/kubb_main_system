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
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type TeacherService interface {
	CreateTeacher(ctx context.Context, input request.TeacherRequestCreate) error
	GetTeacher(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TeacherResponse, *model.PaginationMetadata, error)
	UpdateTeacher(ctx context.Context, uuid string, input request.TeacherRequestUpdate) error
	Toggle(ctx context.Context, uuid string) error
	GetTeacherFilter(ctx context.Context, filter map[string]string) ([]response.TeacherResponseFilter, error)
	CreateTeacherRate(ctx context.Context, input request.TeacherRateRequestCreate, userID int) error
}

type teacherservice struct {
	db *gorm.DB
}

func NewTeacherService() TeacherService {
	return &teacherservice{
		db: config.DB,
	}
}

func (s *teacherservice) CreateTeacherRate(ctx context.Context, input request.TeacherRateRequestCreate, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.TeacherRate{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		TeacherID:       input.TeacherID,
		ClassOfferingID: input.ClassOfferingID,
		HourlyRate:      input.HourlyRate,
		EffectiveDate:   input.EffectiveDate,
		EndDate:         nil,
		Active:          true,
		CreateBy:        userID,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create teacher rate", nil)
		}
		return nil
	})
	return err
}

func (s *teacherservice) CreateTeacher(ctx context.Context, input request.TeacherRequestCreate) error {
	email := strings.TrimSpace(input.Email)
	name := strings.TrimSpace(input.Name)
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
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newteacher).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create teacher", nil)
		}
		newteacher.Code = helper.GenerateCode("T", uint(newteacher.ID))
		if err := tx.Save(&newteacher).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update teacher", nil)
		}
		if len(input.TeacherFacultyRequestCreate) > 0 {
			teacherfaculty := make([]model.TeacherFaculty, 0, len(input.TeacherFacultyRequestCreate))
			for _, t := range input.TeacherFacultyRequestCreate {
				teacherfaculty = append(teacherfaculty, model.TeacherFaculty{
					TeacherID: newteacher.ID,
					FacultyID: t.FacultyID,
				})
			}
			if err := tx.Create(&teacherfaculty).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create teacher facutly", nil)
			}
		}
		return nil
	})
	return err
}

func (s *teacherservice) GetTeacherFilter(
	ctx context.Context,
	filter map[string]string,
) ([]response.TeacherResponseFilter, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.TeacherResponseFilter

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("teachers AS t")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v := filter["faculty_id"]; v != "" {
			facultyIDs := strings.Split(v, ",")

			tx = tx.Where(`
			EXISTS (
				SELECT 1
				FROM teacher_faculty tf
				WHERE tf.teacher_id = t.id
				AND tf.faculty_id IN ?
			)
		`, facultyIDs)
		}

		return tx
	}

	dataQuery := applyFilters(base()).
		Select(`
			t.id AS id,
			t.uuid AS uuid,
			t.code AS code,
			t.email AS email,
			t.name AS name,
			t.date_of_birth AS date_of_birth,
			t.place_of_birth AS place_of_birth,
			t.gender AS gender,
			t.nationality AS nationality,
			t.address AS address,
			t.phone AS phone
		`).
		Order("t.id DESC")

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, fmt.Errorf("fetch teachers: %w", err)
	}

	return data, nil
}

func (s *teacherservice) GetTeacher(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.TeacherResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var data []response.TeacherResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("teachers t")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("t.name LIKE ?", "%"+v+"%")
		}
		return tx
	}

	if err := applyFilters(base()).
		Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count teachers: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).
		Select(`
		t.id AS id,
		t.uuid AS uuid,
		t.code AS code,
		t.email AS email,
		t.name AS name,
		t.date_of_birth AS date_of_birth,
		t.place_of_birth AS place_of_birth,
		t.gender AS gender,
		t.nationality AS nationality,
		t.address AS address,
		t.phone AS phone
	`).Order("t.id DESC").Offset(offset).Limit(pf.PageSize)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch teacher: %w", err)
	}

	teacherIDs := make([]int, 0, len(data))
	for _, teacher := range data {
		teacherIDs = append(teacherIDs, teacher.ID)
	}

	var teacherfaculty []response.TeacherFacultyResponse
	if err := s.db.WithContext(ctx).
		Table("teacher_faculty tf").
		Joins("JOIN faculties f ON f.id = tf.faculty_id").
		Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
		Where("tf.teacher_id IN ?", teacherIDs).Select(`
		tf.teacher_id AS teacher_id,
		f.id AS faculty_id,
		f.code AS faculty_code,
		f.name AS faculty_name,
		p.id AS programme_id,
		p.name AS programme_name
	`).Scan(&teacherfaculty).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch teacher faculty: %w", err)
	}

	teacherFacultyByTeacher := make(
		map[int64][]response.TeacherFacultyResponse,
		len(teacherfaculty),
	)

	for _, tf := range teacherfaculty {
		teacherFacultyByTeacher[int64(tf.TeacherID)] =
			append(teacherFacultyByTeacher[int64(tf.TeacherID)], tf)
	}

	for i := range data {
		data[i].TeacherFacultyResponse =
			teacherFacultyByTeacher[int64(data[i].ID)]

	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *teacherservice) UpdateTeacher(ctx context.Context, uuid string, input request.TeacherRequestUpdate) error {
	if strings.TrimSpace(uuid) == "" {
		return apperror.New(apperror.CodeInvalidInput, "teacher id is required", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var existing model.Teacher
	if err := s.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.CodeNotFound, "teacher not found", nil)
		}
		return helper.MapAcademicError(err, "UPDATE")
	}

	updates := map[string]interface{}{}

	if input.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*input.Email))
		if email == "" {
			return apperror.New(apperror.CodeInvalidInput, "email cannot be empty", nil)
		}
		if !emailRegex.MatchString(email) {
			return apperror.New(apperror.CodeInvalidInput, "email format is invalid", nil)
		}
		if email != existing.Email {
			var dupCount int64
			if err := s.db.WithContext(ctx).
				Table("teachers").
				Where("email = ? AND uuid <> ?", email, uuid).
				Count(&dupCount).Error; err != nil {
				return fmt.Errorf("check existing teacher email: %w", err)
			}
			if dupCount > 0 {
				return apperror.New(apperror.CodeConflict, "a teacher with this email already exists", nil)
			}
		}
		updates["email"] = email
	}

	if input.Name != nil {
		name := strings.ToUpper(strings.TrimSpace(*input.Name))
		if name == "" {
			return apperror.New(apperror.CodeInvalidInput, "name cannot be empty", nil)
		}
		updates["name"] = name
	}

	if input.Dob != nil {
		dob := strings.TrimSpace(*input.Dob)
		if dob == "" {
			return apperror.New(apperror.CodeInvalidInput, "date of birth cannot be empty", nil)
		}
		updates["date_of_birth"] = dob
	}

	if input.Pob != nil {
		pob := strings.ToUpper(strings.TrimSpace(*input.Pob))
		if pob == "" {
			return apperror.New(apperror.CodeInvalidInput, "place of birth cannot be empty", nil)
		}
		updates["place_of_birth"] = pob
	}

	if input.Gender != nil {
		gender := strings.ToUpper(strings.TrimSpace(*input.Gender))
		if gender == "" {
			return apperror.New(apperror.CodeInvalidInput, "gender cannot be empty", nil)
		}
		updates["gender"] = gender
	}

	if input.Nationality != nil {
		nationality := strings.ToUpper(strings.TrimSpace(*input.Nationality))
		if nationality == "" {
			return apperror.New(apperror.CodeInvalidInput, "nationality cannot be empty", nil)
		}
		updates["nationality"] = nationality
	}

	if input.Address != nil {
		address := strings.ToUpper(strings.TrimSpace(*input.Address))
		if address == "" {
			return apperror.New(apperror.CodeInvalidInput, "address cannot be empty", nil)
		}
		updates["address"] = address
	}

	if input.Phone != nil {
		phone := strings.TrimSpace(*input.Phone)
		if phone == "" {
			return apperror.New(apperror.CodeInvalidInput, "phone cannot be empty", nil)
		}
		updates["phone"] = phone
	}

	// Nothing to update at all (no fields, no faculty change) — bail out early.
	if len(updates) == 0 && input.TeacherFacultyRequestCreate == nil {
		return nil
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&existing).Updates(updates).Error; err != nil {
				return helper.MapAcademicError(err, "UPDATE")
			}
		}

		// Only touch faculty links if the caller actually sent a new set.
		if input.TeacherFacultyRequestCreate != nil {
			if err := tx.Where("teacher_id = ?", existing.ID).
				Delete(&model.TeacherFaculty{}).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to clear teacher faculties", nil)
			}

			if len(input.TeacherFacultyRequestCreate) > 0 {
				teacherfaculty := make([]model.TeacherFaculty, 0, len(input.TeacherFacultyRequestCreate))
				for _, t := range input.TeacherFacultyRequestCreate {
					teacherfaculty = append(teacherfaculty, model.TeacherFaculty{
						TeacherID: existing.ID,
						FacultyID: t.FacultyID,
					})
				}
				if err := tx.Create(&teacherfaculty).Error; err != nil {
					return apperror.New(apperror.CodeInternal, "failed to create teacher faculty", nil)
				}
			}
		}

		return nil
	})

	return err
}

func (s *teacherservice) Toggle(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return apperror.New(apperror.CodeInvalidInput, "teacher id is required", nil)
	}
	return utils.ToggleStatus[model.Teacher](ctx, s.db, id)
}
