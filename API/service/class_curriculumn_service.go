package service

import (
	"context"
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

type ClassCurriculumnService interface {
	CreateClassCurriculumn(ctx context.Context, input request.ClassCurriculumnRequestCreate) error
	GetClassCurriculumn(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.ClasCurriculumnResponse, *model.PaginationMetadata, error)
	GetClassCurriculumnWithTeacherRate(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.ClasCurriculumnResponseWithTeacherRate, *model.PaginationMetadata, error)
}

type classcurriculmnservice struct {
	db *gorm.DB
}

func NewClassCurriculumnService() ClassCurriculumnService {
	return &classcurriculmnservice{
		db: config.DB,
	}
}

func (s *classcurriculmnservice) CreateClassCurriculumn(ctx context.Context, input request.ClassCurriculumnRequestCreate) error {
	name := strings.TrimSpace(input.Name)
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newclasscurriculumn := model.ClassCurriculumn{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		Name:    name,
		MajorID: input.MajorID,
		TermID:  input.TermID,
		Active:  true,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newclasscurriculumn).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create class curriculmn", nil)
		}

		if len(input.ClassCurriculumnDetailRequestCreate) > 0 {
			detail := make([]model.ClassCurriculumnDetail, 0, len(input.ClassCurriculumnDetailRequestCreate))
			for _, c := range input.ClassCurriculumnDetailRequestCreate {
				detail = append(detail, model.ClassCurriculumnDetail{
					UUIDBase: base.UUIDBase{
						UUID: helper.GenerateUUID(),
					},
					ClassCurriculumnID: newclasscurriculumn.ID,
					SemesterID:         c.SemesterID,
					StudyYearID:        c.StudyYearID,
					AcademicShiftID:    c.AcademicShiftID,
					MidtermDate:        c.MidtermDate,
					FinalDate:          c.FinalDate,
					TotalStudent:       nil,
					TypeClass:          c.TypeClass,
				})
			}
			if err := tx.Create(&detail).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create dedtail", nil)
			}
		}
		return nil
	})
	return err
}

func (s *classcurriculmnservice) GetClassCurriculumn(
	ctx context.Context,
	pf request.Pagination,
	filter map[string]string,
) ([]response.ClasCurriculumnResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.ClasCurriculumnResponse
	var total int64
	offset := (pf.Page - 1) * pf.PageSize

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("class_curriculums c").
			Joins("LEFT JOIN majors m ON m.id = c.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
			Joins("LEFT JOIN terms t ON t.id = c.term_id").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id")
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
		if v, ok := filter["term_id"]; ok && v != "" {
			tx = tx.Where("t.id = ?", v)
		}
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("g.id = ?", v)
		}
		if v, ok := filter["academic_id"]; ok && v != "" {
			tx = tx.Where("a.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).
		Session(&gorm.Session{}).
		Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count class curriculum: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	dataQuery := applyFilters(base()).Select(`
		c.id AS id,
		c.uuid AS uuid,
		c.name AS name,
		c.active AS active,
		m.id AS major_id,
		m.name AS major_name,
		m.code AS major_code,
		m.duration_period AS major_duration_period,
		m.duration_interval AS major_duration_interval,
		d.id AS department_id,
		d.name AS department_name,
		d.code AS department_code,
		f.id AS faculty_id,
		f.name AS faculty_name,
		p.id AS programme_id,
		p.name AS programme_name,
		t.id AS term_id,
		t.code AS term_code,
		t.name AS term_name,
		g.id AS generation_id,
		g.code AS generation_code,
		g.name AS generation_name,
		a.id AS academic_id,
		a.code AS academic_code,
		a.name AS academic_name
	`).Order("c.id DESC").Offset(offset).Limit(pf.PageSize)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class curriculum: %w", err)
	}

	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	classcurriculumIDs := make([]int, 0, len(data))
	for _, cc := range data {
		classcurriculumIDs = append(classcurriculumIDs, cc.ID)
	}

	var details []response.ClasCurriculumnDetailResponse
	if err := s.db.WithContext(ctx).
		Table("class_curriculum_details cd").
		Joins("LEFT JOIN semesters s ON s.id = cd.semester_id").
		Joins("LEFT JOIN academics a ON a.id = s.academic_id").
		Joins("LEFT JOIN academic_shifts ash ON ash.id = cd.academic_shift_id").
		Where("cd.class_curriculum_id IN ?", classcurriculumIDs).
		Select(`
			cd.id AS id,
			cd.uuid AS uuid,
			cd.class_curriculum_id AS class_curriculum_id,
			s.id AS semester_id,
			s.code AS semester_code,
			s.name AS semester_name,
			a.name AS academic_name,
			cd.study_year_id AS study_year_id,
			ash.id AS academic_shift_id,
			ash.name AS academic_shift_name,
			cd.midterm_date AS midterm_date,
			cd.final_date AS final_date,
			cd.total_student AS total_student,
			cd.type_class AS type_class
		`).Scan(&details).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class curriculum detail: %w", err)
	}

	for i := range details {
		details[i].MidtermDate = helper.FormatDate(details[i].MidtermDate)
		details[i].FinalDate = helper.FormatDate(details[i].FinalDate)
	}

	detailIDs := make([]int, 0, len(details))
	for _, d := range details {
		detailIDs = append(detailIDs, d.ID)
	}

	var classoffer []response.ClassOfferingResponse
	if err := s.db.WithContext(ctx).
		Table("class_offerings co").
		Joins("LEFT JOIN subjects s ON s.id = co.subject_id").
		Where("co.class_curriculum_detail_id IN ?", detailIDs).Select(`
		co.id AS id,
		co.uuid AS uuid,
		co.class_curriculum_detail_id AS class_curriculum_detail_id,
		s.id AS subject_id,
		s.code AS subject_code,
		s.name AS subject_name,
		co.credit AS credit,
		co.passing_score AS passing_score,
		co.total_hour AS total_hour,
		co.status AS status,
		co.total_attendance_for_rexam AS total_attendance_for_rexam,
		co.total_attendance_for_relearn AS total_attendance_for_relearn,
		co.description
	`).Scan(&classoffer).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class curriculum detail: %w", err)
	}

	offerIDs := make([]int, 0, len(classoffer))
	for _, o := range classoffer {
		offerIDs = append(offerIDs, o.ID)
	}
	var student []response.StudentResponseSummary
	if len(offerIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("course_registrations cr").
			Joins("LEFT JOIN student_terms st ON st.id = cr.student_term_id").
			Joins("LEFT JOIN enrollments e ON e.id = st.enrollment_id").
			Joins("LEFT JOIN admissions a ON a.id = e.admission_id").
			Joins("LEFT JOIN students s ON s.id = a.student_id").
			Where("cr.class_offering_id IN ?", offerIDs).
			Select(`
            cr.class_offering_id AS offer_id,
            s.name_kh AS name_kh,
            s.name_en AS name_en,
            s.date_of_birth AS date_of_birth,
            s.gender AS gender,
            s.nationality AS nationality,
            s.phone AS phone,
            s.occupation AS occupation
        `).Scan(&student).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch student: %w", err)
		}
	}

	for i := range student {
		student[i].DateOfBirth = helper.FormatDate(student[i].DateOfBirth)
	}

	studentByOffer := make(map[int][]response.StudentResponseSummary, len(classoffer))
	for _, st := range student {
		studentByOffer[st.OfferID] = append(studentByOffer[st.OfferID], st)
	}

	for i := range classoffer {
		classoffer[i].StudentResponseSummary = studentByOffer[classoffer[i].ID] // add this field if missing
	}

	offeringsByDetail := make(map[int][]response.ClassOfferingResponse, len(details))
	for _, o := range classoffer {
		offeringsByDetail[o.ClassCurriculumnDetailID] = append(offeringsByDetail[o.ClassCurriculumnDetailID], o)
	}
	for i := range details {
		details[i].ClassOfferingResponse = offeringsByDetail[details[i].ID] // adjust field name to match your struct
	}

	detailsByParent := make(map[int][]response.ClasCurriculumnDetailResponse, len(data))
	for _, d := range details {
		detailsByParent[d.ClassCurriculumnID] = append(detailsByParent[d.ClassCurriculumnID], d)
	}
	for i := range data {
		data[i].ClasCurriculumnDetailResponse = detailsByParent[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}

func (s *classcurriculmnservice) GetClassCurriculumnWithTeacherRate(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.ClasCurriculumnResponseWithTeacherRate, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.ClasCurriculumnResponseWithTeacherRate
	var total int64
	offset := (pf.Page - 1) * pf.PageSize

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("class_curriculums c").
			Joins("LEFT JOIN majors m ON m.id = c.major_id").
			Joins("LEFT JOIN departments d ON d.id = m.department_id").
			Joins("LEFT JOIN faculties f ON f.id = d.faculty_id").
			Joins("LEFT JOIN programmes p ON p.id = f.programme_id").
			Joins("LEFT JOIN terms t ON t.id = c.term_id").
			Joins("LEFT JOIN generations g ON g.id = t.generation_id").
			Joins("LEFT JOIN academics a ON a.id = g.academic_id")
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
		if v, ok := filter["term_id"]; ok && v != "" {
			tx = tx.Where("t.id = ?", v)
		}
		if v, ok := filter["generation_id"]; ok && v != "" {
			tx = tx.Where("g.id = ?", v)
		}
		if v, ok := filter["academic_id"]; ok && v != "" {
			tx = tx.Where("a.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).
		Session(&gorm.Session{}).
		Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count class curriculum: %w", err)
	}

	if total == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	dataQuery := applyFilters(base()).Select(`
		c.id AS id,
		c.uuid AS uuid,
		c.name AS name,
		c.active AS active,
		m.id AS major_id,
		m.name AS major_name,
		m.code AS major_code,
		m.duration_period AS major_duration_period,
		m.duration_interval AS major_duration_interval,
		d.id AS department_id,
		d.name AS department_name,
		d.code AS department_code,
		f.id AS faculty_id,
		f.name AS faculty_name,
		p.id AS programme_id,
		p.name AS programme_name,
		t.id AS term_id,
		t.code AS term_code,
		t.name AS term_name,
		g.id AS generation_id,
		g.code AS generation_code,
		g.name AS generation_name,
		a.id AS academic_id,
		a.code AS academic_code,
		a.name AS academic_name
	`).Order("c.id DESC").Offset(offset).Limit(pf.PageSize)

	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class curriculum: %w", err)
	}

	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	classcurriculumIDs := make([]int, 0, len(data))
	for _, cc := range data {
		classcurriculumIDs = append(classcurriculumIDs, cc.ID)
	}

	var details []response.ClasCurriculumnDetailResponseWithTeacherRate
	if err := s.db.WithContext(ctx).
		Table("class_curriculum_details cd").
		Joins("LEFT JOIN semesters s ON s.id = cd.semester_id").
		Joins("LEFT JOIN academics a ON a.id = s.academic_id").
		Joins("LEFT JOIN academic_shifts ash ON ash.id = cd.academic_shift_id").
		Where("cd.class_curriculum_id IN ?", classcurriculumIDs).
		Select(`
			cd.id AS id,
			cd.uuid AS uuid,
			cd.class_curriculum_id AS class_curriculum_id,
			s.id AS semester_id,
			s.code AS semester_code,
			s.name AS semester_name,
			a.name AS academic_name,
			cd.study_year_id AS study_year_id,
			ash.id AS academic_shift_id,
			ash.name AS academic_shift_name,
			cd.midterm_date AS midterm_date,
			cd.final_date AS final_date,
			cd.total_student AS total_student,
			cd.type_class AS type_class
		`).Scan(&details).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class curriculum detail: %w", err)
	}

	for i := range details {
		details[i].MidtermDate = helper.FormatDate(details[i].MidtermDate)
		details[i].FinalDate = helper.FormatDate(details[i].FinalDate)
	}

	detailIDs := make([]int, 0, len(details))
	for _, d := range details {
		detailIDs = append(detailIDs, d.ID)
	}

	var classoffer []response.ClassOfferingResponseWithTeacherRate
	if err := s.db.WithContext(ctx).
		Table("class_offerings co").
		Joins("LEFT JOIN subjects s ON s.id = co.subject_id").
		Joins("LEFT JOIN teacher_rates tr ON tr.class_offer_id = co.id").
		Joins("LEFT JOIN teachers t ON t.id = tr.teacher_id").
		Where("co.class_curriculum_detail_id IN ?", detailIDs).Select(`
		tr.id AS teacher_rate_id,
		co.id AS id,
		co.uuid AS uuid,
		co.class_curriculum_detail_id AS class_curriculum_detail_id,
		s.id AS subject_id,
		s.code AS subject_code,
		s.name AS subject_name,
		co.credit AS credit,
		co.passing_score AS passing_score,
		co.total_hour AS total_hour,
		co.status AS status,
		co.total_attendance_for_rexam AS total_attendance_for_rexam,
		co.total_attendance_for_relearn AS total_attendance_for_relearn,
		co.description,
		t.id AS teacher_id,
		t.name AS teacher_name,
		t.gender AS teacher_gender,
		tr.hourly_rate AS hourly_rate,
		tr.effective_date AS effective_date,
		tr.end_date AS end_date,
		tr.active AS active
	`).Scan(&classoffer).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch class curriculum detail: %w", err)
	}

	for i := range classoffer {
		classoffer[i].EffectiveDate = helper.FormatDate(classoffer[i].EffectiveDate)
	}

	offerIDs := make([]int, 0, len(classoffer))
	for _, o := range classoffer {
		offerIDs = append(offerIDs, o.ID)
	}
	var schedule []response.ScheduleResponse
	if len(offerIDs) > 0 {
		if err := s.db.WithContext(ctx).
			Table("schedules s").
			Joins("LEFT JOIN attendances a ON a.schedule_id = s.id").
			Joins("LEFT JOIN user u ON u.id = s.verify_by").
			Joins("LEFT JOIN school_rooms sr ON sr.id = s.room_id").
			Joins("LEFT JOIN teacher_rates tr ON tr.id = s.teacher_rate_id").
			Joins("LEFT JOIN class_offerings co ON co.id = tr.class_offer_id").
			Where("tr.class_offer_id IN ?", offerIDs).
			Select(`
			u.name_kh AS verify_by,
			s.status AS status,
            s.id AS id,
			s.uuid AS uuid,
			co.id AS class_offering_id,
			s.schedule_date AS schedule_date,
			tr.teacher_id AS teacher_id,
			s.start_time AS start_time,
			s.end_time AS end_time,
			s.total_teach_hours AS total_teach_hours,
			s.description AS description,
			s.active AS active,
			sr.id AS room_id,
			sr.code AS room_code,
			sr.name AS room_name,
			a.id AS attendance_id
        `).Scan(&schedule).Error; err != nil {
			return nil, nil, fmt.Errorf("fetch student: %w", err)
		}
	}

	for i := range schedule {
		schedule[i].SchdeduleDate = helper.FormatDate(schedule[i].SchdeduleDate)
		schedule[i].StartTime = helper.FormatTime(schedule[i].StartTime)
		schedule[i].EndTime = helper.FormatTime(schedule[i].EndTime)
	}

	scheduleByOffer := make(map[int][]response.ScheduleResponse, len(classoffer))
	for _, st := range schedule {
		scheduleByOffer[st.ClassOfferingID] = append(scheduleByOffer[st.ClassOfferingID], st)
	}

	completedHoursByOffer := make(map[int]float64, len(classoffer))
	for _, st := range schedule {
		if st.Status == model.ScheduleStatusCompleted {
			completedHoursByOffer[st.ClassOfferingID] += st.TotalTeachHour
		}
	}

	for i := range classoffer {
		classoffer[i].ScheduleResponse = scheduleByOffer[classoffer[i].ID] // add this field if missing
		classoffer[i].RemainingHour = classoffer[i].TotalHour - completedHoursByOffer[classoffer[i].ID]
	}

	offeringsByDetail := make(map[int][]response.ClassOfferingResponseWithTeacherRate, len(details))
	for _, o := range classoffer {
		offeringsByDetail[o.ClassCurriculumnDetailID] = append(offeringsByDetail[o.ClassCurriculumnDetailID], o)
	}
	for i := range details {
		details[i].ClassOfferingResponseWithTeacherRate = offeringsByDetail[details[i].ID] // adjust field name to match your struct
	}

	detailsByParent := make(map[int][]response.ClasCurriculumnDetailResponseWithTeacherRate, len(data))
	for _, d := range details {
		detailsByParent[d.ClassCurriculumnID] = append(detailsByParent[d.ClassCurriculumnID], d)
	}
	for i := range data {
		data[i].ClasCurriculumnDetailResponseWithTeacherRate = detailsByParent[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
