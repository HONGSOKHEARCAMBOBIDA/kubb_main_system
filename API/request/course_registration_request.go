package request

type CourseRegistrationRequestCreate struct {
	ClassOfferingID           int                         `json:"class_offering_id" gorm:"column:class_offering_id"`
	CourseRegistrationRequest []CourseRegistrationRequest `json:"student_term_id"`
}

type CourseRegistrationRequest struct {
	StudentTermID int `json:"student_term_id" gorm:"column:student_term_id"`
}
