package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	service service.AttendanceService
}

func NewAttendanceController() AttendanceController {
	return AttendanceController{
		service: service.NewAttendanceService(),
	}
}

func (cr *AttendanceController) CreateAttendance(c *gin.Context) {
	var input request.AttendanceRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateAttendance(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *AttendanceController) GetAttendance(c *gin.Context) {
	filter := map[string]string{
		"schedule_id": c.Query("schedule_id"),
	}
	data, err := cr.service.GetAttendance(c, filter)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}

		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
