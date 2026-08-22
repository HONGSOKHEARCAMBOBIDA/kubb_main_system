package request

type ClassCurriculumnDetailRequestCreate struct {
	SemesterID      int    `json:"semester_id" gorm:"column:semester_id"`
	StudyYearID     int    `json:"study_year_id" gorm:"column:study_year_id"`
	AcademicShiftID int    `json:"academic_shift_id" gorm:"column:academic_shift_id"`
	MidtermDate     string `json:"midterm_date" gorm:"column:midterm_date"`
	FinalDate       string `json:"final_date" gorm:"column:final_date"`
	TypeClass       string `json:"type_class" gorm:"column:type_class"`
}
