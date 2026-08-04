package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BuildingController struct {
	service service.BuildingService
}

func NewBuildingController() BuildingController {
	return BuildingController{
		service: service.NewBuildingService(),
	}
}

func (cr *BuildingController) GetBuilding(c *gin.Context) {
	data, err := cr.service.GetBuilding(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
