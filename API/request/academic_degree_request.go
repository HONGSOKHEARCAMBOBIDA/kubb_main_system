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
	AcademicID    *int     `json:"academic_id" validate:"omitempty,gt=0"`
	MajorID       *int     `json:"major_id" validate:"omitempty,gt=0"`
	Name          *string  `json:"name" validate:"omitempty,min=2,max=150"`
	MonthlyFee    *float64 `json:"monthly_fee" validate:"omitempty,gt=0"`
	QuarterlyFee  *float64 `json:"quarterly_fee" validate:"omitempty,gt=0"`
	SemesterlyFee *float64 `json:"semesterly_fee" validate:"omitempty,gt=0"`
	YearlyFee     *float64 `json:"yearly_fee" validate:"omitempty,gt=0"`
	Description   *string  `json:"description" validate:"omitempty,max=1000"`
}
