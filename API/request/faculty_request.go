package request

type FacultyRequestCreate struct {
	ProgrammeID int    `json:"programme_id" gorm:"column:programme_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

type FacultyRequestUpdate struct {
	ProgrammeID *int    `json:"programme_id" gorm:"column:programme_id"`
	Code        *string `json:"code" gorm:"column:code"`
	Name        *string `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
}
