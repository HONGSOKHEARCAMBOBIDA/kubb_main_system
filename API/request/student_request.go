package request

type StudentRequestCreate struct {
	GroupID          int    `json:"group_id" validate:"required"`
	NameKh           string `json:"name_kh" validate:"required,min=2,max=150"`
	NameEn           string `json:"name_en" validate:"required,min=2,max=150"`
	DateOfBirth      string `json:"date_of_birth" validate:"required"`
	Gender           string `json:"gender" validate:"required,min=2,max=30"`
	Nationality      string `json:"nationality" validate:"required,min=2,max=150"`
	Phone            string `json:"phone" validate:"required,min=2,max=150"`
	VillageID        int    `json:"village_id" validate:"omitempty"`
	Occupation       string `json:"occupation" validate:"omitempty,min=2,max=150"`
	AcademicStreamID int    `json:"academic_stream_id" validate:"required"`

	StudentEducationRequestCreate []StudentEducationRequestCreate `json:"student_educations" validate:"omitempty,dive"`
	StudentDocumentRequestCreate  []StudentDocumentRequestCreate  `json:"student_documents" validate:"omitempty,dive"`
	StudentFamilyRequestCreate    []StudentFamilyRequestCreate    `json:"student_family" validate:"omitempty,dive"`
}

type StudentFamilyRequestCreate struct {
	FatherName        string `json:"father_name" validate:"omitempty,min=2,max=150"`
	FatherEnglishName string `json:"father_english_name" validate:"omitempty,min=2,max=150"`
	FatherAge         int    `json:"father_age" validate:"omitempty,gt=0,lte=200"`
	FatherIsAlive     bool   `json:"father_is_alive"`
	FatherPhoneNumber string `json:"father_phone_number" validate:"omitempty,min=2,max=150"`
	FatherOccupation  string `json:"father_occupation" validate:"omitempty,min=2,max=150"`
	FatherWorkplace   string `json:"father_workplace" validate:"omitempty,min=2,max=150"`

	MotherName        string `json:"mother_name" validate:"omitempty,min=2,max=150"`
	MotherEnglishName string `json:"mother_english_name" validate:"omitempty,min=2,max=150"`
	MotherAge         int    `json:"mother_age" validate:"omitempty,gt=0,lte=200"`
	MotherIsAlive     bool   `json:"mother_is_alive"`
	MotherPhoneNumber string `json:"mother_phone_number" validate:"omitempty,min=2,max=150"`
	MotherOccupation  string `json:"mother_occupation" validate:"omitempty,min=2,max=150"`
	MotherWorkplace   string `json:"mother_workplace" validate:"omitempty,min=2,max=150"`
}

type StudentEducationRequestCreate struct {
	Level           string `json:"level" validate:"required,min=2,max=150"`
	SchoolName      string `json:"school_name" validate:"required,min=2,max=150"`
	VillageID       int    `json:"village_id" validate:"omitempty"`
	StartDate       string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate         string `json:"end_date" validate:"required,datetime=2006-01-02,gtefield=StartDate"`
	CertificateDate string `json:"cerificate_date" validate:"omitempty,datetime=2006-01-02"`
	Score           string `json:"score" validate:"omitempty,max=150"`
	Gpa             string `json:"gpa" validate:"omitempty,max=150"`
	Grade           string `json:"grade" validate:"omitempty,max=150"`
}

type StudentDocumentRequestCreate struct {
	DocumentTypeID int    `json:"document_type_id" gorm:"column:document_type_id" validate:"omitempty"`
	RequiredQty    int    `json:"required_qty" gorm:"column:required_qty" validate:"omitempty,gt=0,lte=20"`
	RecieveQty     int    `json:"received_qty" gorm:"column:received_qty" validate:"omitempty,gt=0,lte=20"`
	Remark         string `json:"remark" gorm:"column:remark" validate:"omitempty,min=2,max=150"`
}

type StudentRequestUpdate struct {
	GroupID          int    `json:"group_id" validate:"required"`
	NameKh           string `json:"name_kh" validate:"required,min=2,max=150"`
	NameEn           string `json:"name_en" validate:"required,min=2,max=150"`
	DateOfBirth      string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Gender           string `json:"gender" validate:"required,min=2,max=30"`
	Nationality      string `json:"nationality" validate:"required,min=2,max=150"`
	Phone            string `json:"phone" validate:"required,min=2,max=150"`
	VillageID        int    `json:"village_id" validate:"required"`
	Occupation       string `json:"occupation" validate:"omitempty,min=2,max=150"`
	AcademicStreamID int    `json:"academic_stream_id" validate:"required"`

	StudentEducationRequestUpdate []StudentEducationRequestCreate `json:"student_educations" validate:"omitempty,dive"`
	StudentDocumentRequestUpdate  []StudentDocumentRequestCreate  `json:"student_documents" validate:"omitempty,dive"`
	StudentFamilyRequestUpdate    []StudentFamilyRequestCreate    `json:"student_family" validate:"omitempty,dive"`
}
