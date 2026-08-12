package request

type AcademicDegreeRequestCreate struct {
	AcademicID    int     `json:"academic_id" validate:"required,gt=0"`
	MajorID       int     `json:"major_id" validate:"required,gt=0"`
	Name          string  `json:"name" validate:"required,min=2,max=150"`
	MonthlyFee    float64 `json:"monthly_fee" validate:"required,gt=0"`
	QuarterlyFee  float64 `json:"quarterly_fee" validate:"required,gt=0"`
	SemesterlyFee float64 `json:"semesterly_fee" validate:"required,gt=0"`
	YearlyFee     float64 `json:"yearly_fee" validate:"required,gt=0"`
	Description   string  `json:"description" validate:"max=1000"`
}

type AcademicDegreeRequestUpdate struct {
	AcademicID     *int     `json:"academic_id" gorm:"column:academic_id"`
	MajorID        *int     `json:"major_id" gorm:"column:major_id"`
	Name           *string  `json:"name" gorm:"column:name"`
	MonthlyFee     *float64 `json:"monthly_fee" gorm:"column:monthly_fee"`
	QuarterlyFee   *float64 `json:"quarterly_fee" gorm:"column:quarterly_fee"`
	SemesterlyFeee *float64 `json:"semesterly_fee" gorm:"column:semesterly_fee"`
	YearlyFee      *float64 `json:"yearly_fee" gorm:"column:yearly_fee"`
	Description    *string  `json:"description" gorm:"column:description"`
}
