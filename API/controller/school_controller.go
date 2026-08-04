package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchoolController struct {
	service service.SchoolService
}

func NewSchoolController() SchoolController {
	return SchoolController{
		service: service.NewSchoolService(),
	}
}

func (cr *SchoolController) GetSchool(c *gin.Context) {
	data, err := cr.service.GetSchool(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
