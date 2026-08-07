package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenerationController struct {
	service service.GenerationService
}

func NewGenerationController() GenerationController {
	return GenerationController{
		service: service.NewGenerationService(),
	}
}

func (cr *GenerationController) GetGeneration(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"academic_id": c.Query("academic_id"),
	}

	data, meta, err := cr.service.GetGeneration(c.Request.Context(), request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)
	for i := range data {
		data[i].StartDate = helper.FormatDate(data[i].StartDate)
		data[i].EndDate = helper.FormatDatePtr(data[i].EndDate)
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

func (cr *GenerationController) GetGenerationByAcademic(c *gin.Context) {
	academicID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetGenerationByAcademic(c, academicID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *GenerationController) CreateGeneration(c *gin.Context) {
	var input request.GenerationRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateGeneration(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *GenerationController) UpdateGeneration(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.GenerationRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateGeneration(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *GenerationController) Toggle(c *gin.Context) {
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
