package route

const (

	// Authentication
	Login    = "login"
	Refresh  = "refresh"
	Logout   = "logout"
	Register = "register"

	// User
	UserView   = "user.view"
	UserCreate = "user.create"
	UserUpdate = "user.update/:id"
	UserDelete = "user.delete"

	RoleView             = "role.view"
	RoleCreate           = "role.create"
	RoleUpdate           = "role.update/:id"
	RolePermissionView   = "role.permission.view/:id"
	RolePermissionCreate = "role.permission.create"
	RolePermissionDelete = "role.permission.delete"

	// Author
	AuthorView         = "author.view"
	AuthorCreate       = "author.create"
	AuthorUpdate       = "author.update/:id"
	AuthorToggleStatus = "author.toggle.status/:id"
	AuthorDelete       = "author.delete/:id"

	// Faculty-Department-Program
	// FacultyView    = "faculty.view"
	// DepartmentView = "department.view/:id"

	// Category
	CategoryView         = "category.view"
	CategoryCreate       = "category.create"
	CategoryUpdate       = "category.update/:id"
	CategoryToggleStatus = "category.toggle.status/:id"

	// Cabinet
	CabinetView = "cabinet.view"

	// Filing Cabinet
	FilingCabinetView = "filing.cabinet.view"

	ProgrammesView           = "programmes.view"
	AcademicView             = "academic.view"
	AcademicCreate           = "academic.create"
	AcademicUpdate           = "academic.update/:uuid"
	AcademicToggle           = "academic.toggle/:uuid"
	GenerationView           = "generation.view"
	GenerationViewByAcademic = "generation.view.academic/:id"
	GenerationCreate         = "generation.create"
	GenerationUpdate         = "generation.update/:uuid"
	GenerationToggle         = "generation.toggle/:uuid"
	SemesterView             = "Semester.view"
	SemesterViewByAcademic   = "Semester.view.academic/:id"
	SemesterCreate           = "Semester.create"
	SemesterUpdate           = "Semester.update/:uuid"
	SemesterToggle           = "Semester.toggle/:uuid"
	SchoolView               = "school.view"
	CampuseView              = "campuse.view"
	BuildingView             = "building.view"
	FloorView                = "floor.view"
	SchoolOfficeView         = "school.office.view"
	SchoolRoomView           = "school.room.view"
	// term
	TermView             = "term.view"
	TermViewByGeneration = "term.view.by.generaion/:id"
	TermCreate           = "term.create"
	TermUpdate           = "term.update/:uuid"
	TermToggle           = "term.toggle/:uuid"

	// faculty
	FacultyView            = "faculty.view"
	FacultyViewByProgramme = "faculty.view.by.programmes/:id"
	FacultyCreate          = "faculty.create"
	FacultyUpdate          = "faculty.update/:uuid"
	FacultyToggle          = "faculty.toggle/:uuid"

	// department
	DepartmentView          = "Department.view"
	DepartmentViewByFaculty = "Department.view.by.faculty/:id"
	DepartmentCreate        = "Department.create"
	DepartmentUpdate        = "Department.update/:uuid"
	DepartmentToggle        = "Department.toggle/:uuid"

	// major
	MajorView             = "Major.view"
	MajorViewByDepartment = "Major.view.by.department/:id"
	MajorCreate           = "Major.create"
	MajorUpdate           = "Major.update/:uuid"
	MajorToggle           = "Major.toggle/:uuid"

	// subject
	SubjectView             = "Subject.view"
	SubjectViewByDepartment = "Subject.view.by.major/:id"
	SubjectCreate           = "Subject.create"
	SubjectUpdate           = "Subject.update/:uuid"
	SubjectToggle           = "Subject.toggle/:uuid"

	// major term
	MajorTermView   = "major.term.view"
	MajorTermCreate = "major.term.create"
	MajorTermUpdate = "major.term.update/:uuid"
	MajorTermToggle = "major.term.toggle/:uuid"

	// Academicshift
	AcademicShiftView           = "AcademicShift.view"
	AcademicShiftViewByAcademic = "AcademicShift.view.academic/:id"
	AcademicShiftCreate         = "AcademicShift.create"
	AcademicShiftUpdate         = "AcademicShift.update/:uuid"
	AcademicShiftToggle         = "AcademicShift.toggle/:uuid"

	// AcademicSection
	AcademicSectionView           = "AcademicSection.view"
	AcademicSectionViewByShift    = "AcademicSection.view.by.shift/:id"
	AcademicSectionViewByAcademic = "AcademicSection.view.academic/:id"
	AcademicSectionCreate         = "AcademicSection.create"
	AcademicSectionUpdate         = "AcademicSection.update/:uuid"
	AcademicSectionToggle         = "AcademicSection.toggle/:uuid"

	// AcademicDegree
	AcademicDegreeView           = "AcademicDegree.view"
	AcademicDegreeViewByAcademic = "AcademicDegree.view.by.academic/:id"
	AcademicDegreeCreate         = "AcademicDegree.create"
	AcademicDegreeUpdate         = "AcademicDegree.update/:uuid"
	AcademicDegreeToggle         = "AcademicDegree.toggle/:uuid"

	// FeediscountGroup
	FeediscountGroupView   = "FeediscountGroup.view"
	FeediscountGroupCreate = "FeediscountGroup.create"
	FeediscountGroupUpdate = "FeediscountGroup.update/:uuid"
	FeediscountGroupToggle = "FeediscountGroup.toggle/:uuid"

	// Documenttype
	DocumentTypeView = "document.type.view"

	// Location
	ProvinceView = "province.view"
	DistrictView = "district.view/:id"
	CommunceView = "communce.view/:id"
	VillageView  = "village.view/:id"

	// Academic Stream
	AcademicStreamView = "academic.stream.view"

	// Student
	StudentCreate = "student.create"
	StudentView   = "student.view"
	StudentUpdate = "student.update/:id"

	// SchoolarshipGroup
	SchoolarshipGroupView   = "SchoolarshipGroup.view"
	SchoolarshipGroupCreate = "SchoolarshipGroup.create"
	SchoolarshipGroupUpdate = "SchoolarshipGroup.update/:uuid"
	SchoolarshipGroupToggle = "SchoolarshipGroup.toggle/:uuid"
)
