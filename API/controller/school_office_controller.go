package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchoolOfficeController struct {
	service service.SchoolOfficeService
}

func NewSchoolOfficeController() SchoolOfficeController {
	return SchoolOfficeController{
		service: service.NewSchoolOfficeService(),
	}
}

func (cr *SchoolOfficeController) GetSchoolOffice(c *gin.Context) {
	data, err := cr.service.GetSchoolOffice(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
