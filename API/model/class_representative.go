package model

import "mysql/model/base"

type ClassRepresentative struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnDetailID int                         `json:"class_curriculum_detail_id" gorm:"column:class_curriculum_detail_id"`
	StudentTermID            int                         `json:"student_term_id" gorm:"column:student_term_id"`
	Position                 ClassRepresentativePosition `gorm:"type:enum('president','vice_president_1','vice_president_2'); not null"`
	StartDate                string                      `json:"start_date" gorm:"column:start_date"`
	EndDate                  string                      `json:"end_date" gorm:"column:end_date"`
	Isactive                 bool                        `json:"is_active" gorm:"column:is_active"`
}

func (ClassRepresentative) TableName() string {
	return "class_representatives"
}
