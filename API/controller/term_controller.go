package controller

import (
	"context"
	"errors"
	"net/http"

	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"

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

	for i := range data {
		data[i].StartDate = helper.FormatDate(data[i].StartDate)
		data[i].EndDate = helper.FormatDate(data[i].EndDate)
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

func (cr *TermController) CreateTerm(c *gin.Context) {
	var input request.TermRequestCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateTerm(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *TermController) UpdateTerm(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.TermRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateTerm(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *TermController) Toggle(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		share.ResponseError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := cr.service.Toggle(c, id); err != nil {
		share.RespondServiceError(c, err)
		return
	}
}
