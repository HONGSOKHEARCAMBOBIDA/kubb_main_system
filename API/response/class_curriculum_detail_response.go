package response

import "mysql/model/base"

type ClasCurriculumnDetailResponse struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnID          int                           `json:"class_curriculum_id" gorm:"column:class_curriculum_id"`
	SemesterID                  int                           `json:"semester_id" gorm:"column:semester_id"`
	SemesterCode                string                        `json:"semester_code"`
	SemesterName                string                        `json:"semester_name"`
	AcademicID                  int                           `json:"academic_id" gorm:"column:academic_id"`
	AcademicName                string                        `json:"academic_name"`
	AcademicCode                string                        `json:"academic_code"`
	StudyYearID                 int                           `json:"study_year_id" gorm:"column:study_year_id"`
	AcademicShiftID             int                           `json:"academic_shift_id" gorm:"column:academic_shift_id"`
	AcademicShiftName           string                        `json:"academic_shift_name"`
	StartDate                   string                        `json:"start_date" gorm:"column:start_date"`
	MidtermDate                 string                        `json:"midterm_date" gorm:"column:midterm_date"`
	FinalDate                   string                        `json:"final_date" gorm:"column:final_date"`
	TotalStudent                *int                          `json:"total_student" gorm:"column:total_student"`
	TypeClass                   string                        `json:"type_class" gorm:"column:type_class"`
	ClassOfferingResponse       []ClassOfferingResponse       `json:"class_offering" gorm:"-"`
	ClassRepresentativeResponse []ClassRepresentativeResponse `json:"class_representative" gorm:"-"`
}

type ClasCurriculumnDetailResponseWithTeacherRate struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnID                   int                                    `json:"class_curriculum_id" gorm:"column:class_curriculum_id"`
	SemesterID                           int                                    `json:"semester_id" gorm:"column:semester_id"`
	SemesterCode                         string                                 `json:"semester_code"`
	SemesterName                         string                                 `json:"semester_name"`
	AcademicID                           int                                    `json:"academic_id" gorm:"column:academic_id"`
	AcademicName                         string                                 `json:"academic_name"`
	AcademicCode                         string                                 `json:"academic_code"`
	StudyYearID                          int                                    `json:"study_year_id" gorm:"column:study_year_id"`
	AcademicShiftID                      int                                    `json:"academic_shift_id" gorm:"column:academic_shift_id"`
	AcademicShiftName                    string                                 `json:"academic_shift_name"`
	StartDate                            string                                 `json:"start_date" gorm:"column:start_date"`
	MidtermDate                          string                                 `json:"midterm_date" gorm:"column:midterm_date"`
	FinalDate                            string                                 `json:"final_date" gorm:"column:final_date"`
	TotalStudent                         *int                                   `json:"total_student" gorm:"column:total_student"`
	TypeClass                            string                                 `json:"type_class" gorm:"column:type_class"`
	ClassOfferingResponseWithTeacherRate []ClassOfferingResponseWithTeacherRate `json:"class_offering" gorm:"-"`
}
