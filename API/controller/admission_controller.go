package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdmissionController struct {
	service service.AdmissionService
}

func NewAdmissionController() AdmissionController {
	return AdmissionController{
		service: service.NewAdmissionService(),
	}
}

func (cr *AdmissionController) GetAdmission(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"name": c.Query("name"),
	}

	data, meta, err := cr.service.GetAdmission(c.Request.Context(), request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)
	for i := range data {
		data[i].Date = helper.FormatDate(data[i].Date)
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponsePagination(c, 200, data, meta)
}
