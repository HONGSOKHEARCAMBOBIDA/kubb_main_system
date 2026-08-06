package controller

import (
	"context"
	"errors"
	"net/http"

	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"

	"github.com/gin-gonic/gin"
)

type TermController struct {
	service service.TermService
}

func NewTermController() TermController {
	return TermController{
		service: service.NewTermService(),
	}
}

func (cr *TermController) GetTerm(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"generation_id": c.Query("generation_id"),
		"academic_id":   c.Query("academic_id"),
	}

	data, meta, err := cr.service.GetTerm(c.Request.Context(), request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)

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
