package request

type AcademicSectionRequestCreate struct {
	MajorID     int    `json:"major_id" gorm:"column:major_id"`
	ShiftID     int    `json:"shift_id" gorm:"column:shift_id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

type AcademicSectionRequestUpdate struct {
	MajorID     *int    `json:"major_id" gorm:"column:major_id"`
	ShiftID     *int    `json:"shift_id" gorm:"column:shift_id"`
	Name        *string `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
}
