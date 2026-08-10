package request

type MajorTermReqeustCreate struct {
	MajorID []int `json:"major_id" gorm:"column:major_id"`
	TermID  int   `json:"term_id" gorm:"column:term_id"`
}

type MajorTermReqeustUpdate struct {
	MajorID int `json:"major_id" gorm:"column:major_id"`
	TermID  int `json:"term_id" gorm:"column:term_id"`
}
