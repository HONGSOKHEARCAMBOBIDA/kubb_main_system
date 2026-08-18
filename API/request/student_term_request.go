package request

type StudentTermRequestCreate struct {
	SemesterID  int `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID int `json:"study_year_id" gorm:"column:study_year_id"`
}
