package request

type AcademicShiftRequestCreate struct {
	AcademicID  int    `json:"academic_id" gorm:"column:academic_id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

type AcademicShiftRequestUpdate struct {
	AcademicID  *int    `json:"academic_id" gorm:"column:academic_id"`
	Name        *string `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
}
