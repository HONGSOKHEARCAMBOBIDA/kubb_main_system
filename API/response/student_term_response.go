package response

import "mysql/model/base"

type StudentTermResponse struct {
	base.ModelBase
	base.UUIDBase
	EnrollmentID int    `json:"enrollment_id" gorm:"column:enrollment_id"`
	SemesterID   int    `json:"semester_id"`
	SemesterCode string `json:"semester_code"`
	SemesterName string `json:"semester_name"`
	AcademicID   int    `json:"academic_id" gorm:"column:academic_id"`
	AcademicName string `json:"academic_name"`
	//	FeeResponse  []FeeResponse `json:"fee" gorm:"-"`
	StudyYearID       int                 `json:"study_year_id" gorm:"column:study_year_id"`
	Active            bool                `json:"active" gorm:"column:active"`
	Status            string              `json:"status" gorm:"column:status"`
	GpaRecordResponse []GpaRecordResponse `json:"gpa_record" gorm:"-"`
}

type StudentTermResponsebyFilter struct {
	base.ModelBase
	base.UUIDBase
	StudentID     int    `json:"student_id"`
	StudentNameKh string `json:"student_name_kh"`
	StudentNameEn string `json:"student_name_en"`
	StudentGender string `json:"student_gender"`
	SemesterID    int    `json:"semester_id"`
	SemesterName  string `json:"semester_name"`
	StudyYearID   int    `json:"study_year_id"`
	TermID        int    `json:"term_id"`
	TermCode      string `json:"term_code"`
	TermName      string `json:"term_name"`
	MajorID       int    `json:"major_id"`
	MajorCode     string `json:"major_code"`
	MajorName     string `json:"major_name"`
	ProgrammID    int    `json:"programm_id"`
	ProgrammName  string `json:"programm_name" gorm:"column:programm_name"`
}
