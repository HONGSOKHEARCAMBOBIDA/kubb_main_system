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

	"gorm.io/gorm"
)

type AttendanceService interface {
	CreateAttendance(ctx context.Context, input request.AttendanceRequestCreate) error
	GetAttendance(ctx context.Context, filter map[string]string) ([]response.AttendanceResponse, error)
}

type attendanceservice struct {
	db *gorm.DB
}

func NewAttendanceService() AttendanceService {
	return &attendanceservice{
		db: config.DB,
	}
}

func (s *attendanceservice) CreateAttendance(ctx context.Context, input request.AttendanceRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	newdata := model.Attendance{
		UUIDBase: base.UUIDBase{
			UUID: helper.GenerateUUID(),
		},
		ScheduleID:     input.ScheduleID,
		AttendanceDate: input.AttendanceDate,
		Topic:          input.Topic,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create attendance", nil)
		}
		if len(input.AttendanceDetailRequestCreate) > 0 {
			detail := make([]model.AttendanceDetail, 0, len(input.AttendanceDetailRequestCreate))
			for _, d := range input.AttendanceDetailRequestCreate {
				detail = append(detail, model.AttendanceDetail{
					UUIDBase: base.UUIDBase{
						UUID: helper.GenerateUUID(),
					},
					AttendanceID:          newdata.ID,
					CourseRegistrationsID: d.CourseRegistrationsID,
					Status:                d.Status,
					Note:                  d.Note,
				})
			}
			if err := tx.Create(&detail).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create attendance detail", nil)
			}
		}
		return nil
	})
	return err
}

func (s *attendanceservice) GetAttendance(ctx context.Context, filter map[string]string) ([]response.AttendanceResponse, error) {
	var data []response.AttendanceResponse

	baseQuery := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("attendances a").
			Joins("JOIN schedules s ON s.id = a.schedule_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["schedule_id"]; ok && v != "" {
			tx = tx.Where("s.id = ?", v)
		}
		return tx
	}

	dataQuery := applyFilters(baseQuery()).
		Select(`
			a.id AS id,
			a.uuid AS uuid,
			a.attendance_date AS attendance_date,
			a.topic AS topic
		`)
	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, fmt.Errorf("fetch attendance: %w", err)
	}

	for i := range data {
		data[i].AttendanceDate = helper.FormatDate(data[i].AttendanceDate)
	}

	if len(data) == 0 {
		return data, nil
	}

	attendanceIDs := make([]int, 0, len(data))
	for _, attendance := range data {
		attendanceIDs = append(attendanceIDs, attendance.ID)
	}

	var detail []response.AttendanceDetailResponse
	if err := s.db.WithContext(ctx).Table("attendance_details ad").
		Joins("LEFT JOIN course_registrations cr ON cr.id = ad.course_registration_id").
		Joins("LEFT JOIN student_terms st ON st.id = cr.student_term_id").
		Joins("LEFT JOIN enrollments e ON e.id = st.enrollment_id").
		Joins("LEFT JOIN admissions adm ON adm.id = e.admission_id").
		Joins("LEFT JOIN students stu ON stu.id = adm.student_id").
		Where("ad.attendance_id IN ?", attendanceIDs).
		Select(`
			ad.status AS status,
			ad.id AS id,
			ad.uuid AS uuid,
			ad.attendance_id AS attendance_id,
			ad.course_registration_id AS course_registration_id,
			stu.name_kh AS name_kh,
			stu.name_en AS name_en,
			stu.date_of_birth AS date_of_birth,
			stu.gender AS gender,
			stu.phone AS phone
		`).Scan(&detail).Error; err != nil {
		return nil, fmt.Errorf("fetch detail: %w", err)
	}

	for i := range detail {
		detail[i].Dob = helper.FormatDate(detail[i].Dob)
	}

	// group details by attendance_id
	detailsByAttendance := make(map[int][]response.AttendanceDetailResponse, len(data))
	for _, d := range detail {
		detailsByAttendance[d.AttendanceID] = append(detailsByAttendance[d.AttendanceID], d)
	}

	for i := range data {
		data[i].AttendanceDetailResponse = detailsByAttendance[data[i].ID]
	}

	return data, nil
}
