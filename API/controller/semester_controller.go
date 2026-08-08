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

type SemesterController struct {
	service service.SemesterService
}

func NewSemesterController() SemesterController {
	return SemesterController{
		service: service.NewSemesterService(),
	}
}

func (cr *SemesterController) GetSemester(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"academic_id": c.Query("academic_id"),
	}

	data, meta, err := cr.service.GetSemester(c.Request.Context(), request.Pagination{
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

func (cr *SemesterController) GetSemesterByAcademic(c *gin.Context) {
	academicID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetSemesterByAcademic(c, academicID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *SemesterController) CreateSemester(c *gin.Context) {
	var input request.SemesterRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateSemester(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *SemesterController) UpdateSemester(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.SemesterRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateSemester(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *SemesterController) Toggle(c *gin.Context) {
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
