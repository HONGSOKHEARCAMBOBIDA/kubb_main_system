package request

type TeacherRequestCreate struct {
	Email       string `json:"email" gorm:"column:email"`
	Name        string `json:"name" gorm:"column:name"`
	Dob         string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Pob         string `json:"place_of_birth" gorm:"column:place_of_birth"`
	Gender      string `json:"gender" gorm:"column:gender"`
	Nationality string `json:"nationality" gorm:"column:nationality"`
	Address     string `json:"address" gorm:"column:address"`
	Phone       string `json:"phone" gorm:"column:phone"`
	FacultyID   int    `json:"faculty_id" gorm:"column:faculty_id"`
}

type TeacherRequestUpdate struct {
	Email       string `json:"email" gorm:"column:email"`
	Name        string `json:"name" gorm:"column:name"`
	Dob         string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Pob         string `json:"place_of_birth" gorm:"column:place_of_birth"`
	Gender      string `json:"gender" gorm:"column:gender"`
	Nationality string `json:"nationality" gorm:"column:nationality"`
	Address     string `json:"address" gorm:"column:address"`
	Phone       string `json:"phone" gorm:"column:phone"`
	FacultyID   int    `json:"faculty_id" gorm:"column:faculty_id"`
}
