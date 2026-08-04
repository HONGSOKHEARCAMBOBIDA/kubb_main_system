package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchoolRoomController struct {
	service service.SchoolRoomService
}

func NewSchoolRoomController() SchoolRoomController {
	return SchoolRoomController{
		service: service.NewSchoolRoomService(),
	}
}

func (cr *SchoolRoomController) GetSchoolRoom(c *gin.Context) {
	data, err := cr.service.GetSchoolRoom(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
