package request

type TeacherFacultyRequestCreate struct {
	FacultyID int `json:"faculty_id" gorm:"column:faculty_id"`
}
