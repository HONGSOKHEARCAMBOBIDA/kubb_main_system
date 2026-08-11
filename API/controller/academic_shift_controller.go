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

type AcademicShiftController struct {
	service service.AcademicShiftService
}

func NewAcademicShiftController() AcademicShiftController {
	return AcademicShiftController{
		service: service.NewAcademicShiftService(),
	}
}

func (cr *AcademicShiftController) GetAcademicShift(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"academic_id": c.Query("academic_id"),
	}

	data, meta, err := cr.service.GetAcademicShift(c.Request.Context(), request.Pagination{
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

func (cr *AcademicShiftController) GetAcademicShiftByAcademic(c *gin.Context) {
	academicID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetAcademicShiftByAcademic(c, academicID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AcademicShiftController) CreateAcademicShift(c *gin.Context) {
	var input request.AcademicShiftRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateAcademicShift(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *AcademicShiftController) UpdateAcademicShift(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.AcademicShiftRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateAcademicShift(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *AcademicShiftController) Toggle(c *gin.Context) {
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
