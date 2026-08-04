package request

type GenerationRequestCreate struct {
	AcademicID  int    `json:"academic_id" gorm:"column:academic_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	StartDate   string `json:"start_date" gorm:"column:start_date"`
	Description string `json:"description" gorm:"column:description"`
}

type GenerationRequestUpdate struct {
	AcademicID  *int    `json:"academic_id" gorm:"column:academic_id"`
	Code        *string `json:"code" gorm:"column:code"`
	Name        *string `json:"name" gorm:"column:name"`
	StartDate   *string `json:"start_date" gorm:"column:start_date"`
	EndDate     *string `json:"end_date" gorm:"column:end_date"`
	Description *string `json:"description" gorm:"column:description"`
}
