package request

type TeacherRateRequestCreate struct {
	TeacherID       int     `json:"teacher_id" gorm:"column:teacher_id"`
	ClassOfferingID int     `json:"class_offer_id" gorm:"column:class_offer_id"`
	HourlyRate      float64 `json:"hourly_rate" gorm:"column:hourly_rate"`
	EffectiveDate   string  `json:"effective_date" gorm:"column:effective_date"`
}
