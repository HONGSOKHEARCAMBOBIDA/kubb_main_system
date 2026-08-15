package request

type AcademicSectionRequestCreate struct {
	MajorID     int    `json:"major_id" gorm:"column:major_id"`
	ShiftID     int    `json:"shift_id" gorm:"column:shift_id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Type        int    `json:"type" gorm:"column:type"`
}

type AcademicSectionRequestUpdate struct {
	MajorID     *int    `json:"major_id" gorm:"column:major_id"`
	ShiftID     *int    `json:"shift_id" gorm:"column:shift_id"`
	Name        *string `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
	Type        *int    `json:"type" gorm:"column:type"`
}
