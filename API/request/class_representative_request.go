package request

type ClassRepresentativeRequestCreate struct {
	ClassCurriculumnDetailID   int                          `json:"class_curriculum_detail_id" gorm:"column:class_curriculum_detail_id"`
	ClassRepresentativeRequest []ClassRepresentativeRequest `json:"student_term_id"`
}

type ClassRepresentativeRequest struct {
	StudentTermID int    `json:"student_term_id" gorm:"column:student_term_id"`
	StartDate     string `json:"start_date" gorm:"column:start_date"`
	EndDate       string `json:"end_date" gorm:"column:end_date"`
}
