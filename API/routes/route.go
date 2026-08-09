package routes

import (
	"mysql/constant/permission"
	"mysql/constant/route"
	"mysql/controller"
	"mysql/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authcontroller := controller.NewAuthController()
	rolecontroller := controller.NewRoleHasPermissionController()
	programmescontroller := controller.NewProgrammescontroller()
	academiccontroller := controller.NewAcademicController()
	generationcontroller := controller.NewGenerationController()
	schoolcontroller := controller.NewSchoolController()
	campusecontroller := controller.NewCampuseController()
	buildingcontroller := controller.NewBuildingController()
	floorcontroller := controller.NewFloorController()
	schoolofficecontroller := controller.NewSchoolOfficeController()
	schoolroomcontroller := controller.NewSchoolRoomController()
	termcontroller := controller.NewTermController()
	semestercontroller := controller.NewSemesterController()
	facultycontroller := controller.NewFacultyController()
	public := r.Group("/api/v1")
	public.Use(middleware.APIKeyAuth())
	{
		public.POST(route.Login, authcontroller.Login)
		public.POST(route.Refresh, authcontroller.Refresh)
	}
	auth := r.Group("/api/v1")
	auth.Use(middleware.APIKeyAuth())
	auth.Use(middleware.AuthMiddleware())
	{
		// User
		auth.POST(route.Register, middleware.PermissionMiddleware(permission.UserCreate), authcontroller.Register)
		auth.GET(route.UserView, middleware.PermissionMiddleware(permission.UserView), authcontroller.GetUser)
		auth.PUT(route.UserUpdate, middleware.PermissionMiddleware(permission.UserUpdate), authcontroller.Update)

		// Role
		auth.GET(route.RoleView, middleware.PermissionMiddleware(permission.CRUDPERMISSION), rolecontroller.GetRole)
		auth.PUT(route.RoleUpdate, middleware.PermissionMiddleware(permission.CRUDPERMISSION), rolecontroller.UpdateRole)
		auth.GET(route.RolePermissionView, middleware.PermissionMiddleware(permission.CRUDPERMISSION), rolecontroller.GetRolePermission)
		auth.POST(route.RolePermissionCreate, middleware.PermissionMiddleware(permission.CRUDPERMISSION), rolecontroller.CreateRoleHasPermission)
		auth.DELETE(route.RolePermissionDelete, middleware.PermissionMiddleware(permission.CRUDPERMISSION), rolecontroller.DeleteRoleHasPermission)

		// Programmes
		auth.GET(route.ProgrammesView, middleware.PermissionMiddleware(permission.ProgrammesView), programmescontroller.GetProgrammes)

		// Academic
		auth.GET(route.AcademicView, middleware.PermissionMiddleware(permission.CRUDACADEMIC), academiccontroller.GetAcademic)
		auth.POST(route.AcademicCreate, middleware.PermissionMiddleware(permission.CRUDACADEMIC), academiccontroller.CreateAcademic)
		auth.PUT(route.AcademicUpdate, middleware.PermissionMiddleware(permission.CRUDACADEMIC), academiccontroller.UpdateAcademic)
		auth.PUT(route.AcademicToggle, middleware.PermissionMiddleware(permission.CRUDACADEMIC), academiccontroller.Toggle)
		// Generation
		auth.GET(route.GenerationView, middleware.PermissionMiddleware(permission.CRUDGENERATION), generationcontroller.GetGeneration)
		auth.GET(route.GenerationViewByAcademic, middleware.PermissionMiddleware(permission.CRUDGENERATION), generationcontroller.GetGenerationByAcademic)
		auth.POST(route.GenerationCreate, middleware.PermissionMiddleware(permission.CRUDGENERATION), generationcontroller.CreateGeneration)
		auth.PUT(route.GenerationUpdate, middleware.PermissionMiddleware(permission.CRUDGENERATION), generationcontroller.UpdateGeneration)
		auth.PUT(route.GenerationToggle, middleware.PermissionMiddleware(permission.CRUDGENERATION), generationcontroller.Toggle)

		// Semester
		auth.GET(route.SemesterView, middleware.PermissionMiddleware(permission.CRUDSEMESTER), semestercontroller.GetSemester)
		auth.GET(route.SemesterViewByAcademic, middleware.PermissionMiddleware(permission.CRUDSEMESTER), semestercontroller.GetSemesterByAcademic)
		auth.POST(route.SemesterCreate, middleware.PermissionMiddleware(permission.CRUDSEMESTER), semestercontroller.CreateSemester)
		auth.PUT(route.SemesterUpdate, middleware.PermissionMiddleware(permission.CRUDSEMESTER), semestercontroller.UpdateSemester)
		auth.PUT(route.SemesterToggle, middleware.PermissionMiddleware(permission.CRUDSEMESTER), semestercontroller.Toggle)
		// School
		auth.GET(route.SchoolView, middleware.PermissionMiddleware(permission.CRUDSCHOOL), schoolcontroller.GetSchool)

		// Campuse
		auth.GET(route.CampuseView, middleware.PermissionMiddleware(permission.CRUDECAMPUSE), campusecontroller.GetCampuse)

		// Building
		auth.GET(route.BuildingView, middleware.PermissionMiddleware(permission.CRUDEBUILDING), buildingcontroller.GetBuilding)

		// Floor
		auth.GET(route.FloorView, middleware.PermissionMiddleware(permission.CRUDFLOOR), floorcontroller.GetFloor)

		// SchoolOffice
		auth.GET(route.SchoolOfficeView, middleware.PermissionMiddleware(permission.CRUDSCHOOLOFFICE), schoolofficecontroller.GetSchoolOffice)

		// SchoolRoom
		auth.GET(route.SchoolRoomView, middleware.PermissionMiddleware(permission.CRUDSCHOOLROOM), schoolroomcontroller.GetSchoolRoom)

		// Term
		auth.GET(route.TermView, middleware.PermissionMiddleware(permission.CRUDTERM), termcontroller.GetTerm)
		auth.POST(route.TermCreate, middleware.PermissionMiddleware(permission.CRUDTERM), termcontroller.CreateTerm)
		auth.PUT(route.TermUpdate, middleware.PermissionMiddleware(permission.CRUDTERM), termcontroller.UpdateTerm)
		auth.PUT(route.TermToggle, middleware.PermissionMiddleware(permission.CRUDTERM), termcontroller.Toggle)

		// Faculty
		auth.GET(route.FacultyView, middleware.PermissionMiddleware(permission.CRUDFACULTY), facultycontroller.GetFaculty)
		auth.GET(route.FacultyViewByProgramme, middleware.PermissionMiddleware(permission.CRUDFACULTY), facultycontroller.GetFacultyByProgrammes)
		auth.POST(route.FacultyCreate, middleware.PermissionMiddleware(permission.CRUDFACULTY), facultycontroller.CreateFaculty)
		auth.PUT(route.FacultyUpdate, middleware.PermissionMiddleware(permission.CRUDFACULTY), facultycontroller.UpdateFaculty)
		auth.PUT(route.FacultyToggle, middleware.PermissionMiddleware(permission.CRUDFACULTY), facultycontroller.Toggle)

	}
}
