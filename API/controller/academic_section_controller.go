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

type AcademicSectionController struct {
	service service.AcademicSectionService
}

func NewAcademicSectionController() AcademicSectionController {
	return AcademicSectionController{
		service: service.NewAcademicSectionService(),
	}
}

func (cr *AcademicSectionController) GetAcademicSectionByShift(c *gin.Context) {
	departmentID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetAcademicSectionByShift(c, departmentID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AcademicSectionController) GetAcademicSection(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"academic_id":   c.Query("academic_id"),
		"shift_id":      c.Query("shift_id"),
		"programme_id":  c.Query("programme_id"),
		"faculty_id":    c.Query("faculty_id"),
		"department_id": c.Query("department_id"),
		"major_id":      c.Query("major_id"),
	}

	data, meta, err := cr.service.GetAcademicSection(c.Request.Context(), request.Pagination{
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

func (cr *AcademicSectionController) CreateAcademicSection(c *gin.Context) {
	var input request.AcademicSectionRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateAcademicSection(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *AcademicSectionController) UpdateAcademicSection(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.AcademicSectionRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateAcademicSection(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *AcademicSectionController) Toggle(c *gin.Context) {
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
