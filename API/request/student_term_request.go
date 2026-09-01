package request

type StudentTermRequestCreate struct {
	SemesterID  int `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID int `json:"study_year_id" gorm:"column:study_year_id"`
}

type StudentTermRequestv2 struct {
	EnrollmentID int `json:"enrollment_id"`
	SemesterID   int `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID  int `json:"study_year_id" gorm:"column:study_year_id"`
}

type StudentTermRequestUpdate struct {
	SemesterID  int `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID int `json:"study_year_id" gorm:"column:study_year_id"`
}
