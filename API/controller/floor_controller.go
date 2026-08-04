package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FloorController struct {
	service service.FloorService
}

func NewFloorController() FloorController {
	return FloorController{
		service: service.NewFloorService(),
	}
}

func (cr *FloorController) GetFloor(c *gin.Context) {
	data, err := cr.service.GetFloor(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
