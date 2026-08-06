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
	FacultyView    = "faculty.view"
	DepartmentView = "department.view/:id"

	// Category
	CategoryView         = "category.view"
	CategoryCreate       = "category.create"
	CategoryUpdate       = "category.update/:id"
	CategoryToggleStatus = "category.toggle.status/:id"

	// Cabinet
	CabinetView = "cabinet.view"

	// Filing Cabinet
	FilingCabinetView = "filing.cabinet.view"

	ProgrammesView   = "programmes.view"
	AcademicView     = "academic.view"
	AcademicCreate   = "academic.create"
	AcademicUpdate   = "academic.update/:uuid"
	AcademicToggle   = "academic.toggle/:uuid"
	GenerationView   = "generation.view"
	GenerationCreate = "generation.create"
	GenerationUpdate = "generation.update/:uuid"
	GenerationToggle = "generation.toggle/:uuid"
	SchoolView       = "school.view"
	CampuseView      = "campuse.view"
	BuildingView     = "building.view"
	FloorView        = "floor.view"
	SchoolOfficeView = "school.office.view"
	SchoolRoomView   = "school.room.view"
	TermView         = "term.view"
)
