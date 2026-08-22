package request

type ClassCurriculumnRequestCreate struct {
	Name                                string                                `json:"name" gorm:"column:name"`
	MajorID                             int                                   `json:"major_id" gorm:"column:major_id"`
	TermID                              int                                   `json:"term_id" gorm:"column:term_id"`
	ClassCurriculumnDetailRequestCreate []ClassCurriculumnDetailRequestCreate `json:"class_curriclumn_details"`
}
