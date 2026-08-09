package request

type DepartmentRequestCreate struct {
	FacultyID   int    `json:"faculty_id" gorm:"column:faculty_id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

type DepartmentRequestUpdate struct {
	FacultyID   *int    `json:"faculty_id" gorm:"column:faculty_id"`
	Code        *string `json:"code" gorm:"column:code"`
	Name        *string `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
}
