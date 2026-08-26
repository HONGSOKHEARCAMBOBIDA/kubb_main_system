package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	service service.ScheduleService
}

func NewScheduleController() ScheduleController {
	return ScheduleController{
		service: service.NewScheduleService(),
	}
}

func (cr *ScheduleController) CreateSchedule(c *gin.Context) {
	var input request.ScheduleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateScheduleService(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}
