package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CampuseController struct {
	service service.CampuseService
}

func NewCampuseController() CampuseController {
	return CampuseController{
		service: service.NewCampuseService(),
	}
}

func (cr *CampuseController) GetCampuse(c *gin.Context) {
	data, err := cr.service.GetCampuse(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
