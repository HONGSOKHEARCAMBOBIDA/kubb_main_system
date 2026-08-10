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

type MajorTermController struct {
	service service.MajorTermService
}

func NewMajorTermController() MajorTermController {
	return MajorTermController{
		service: service.NewMajorTermService(),
	}
}

func (cr *MajorTermController) GetMajorTerm(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"programme_id":  c.Query("programme_id"),
		"faculty_id":    c.Query("faculty_id"),
		"department_id": c.Query("department_id"),
		"academic_id":   c.Query("academic_id"),
		"generation_id": c.Query("generation_id"),
		"term_id":       c.Query("term_id"),
	}

	data, meta, err := cr.service.GetMajorTerm(c.Request.Context(), request.Pagination{
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

func (cr *MajorTermController) CreateMajorTerm(c *gin.Context) {
	var input request.MajorTermReqeustCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateMajorTerm(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *MajorTermController) UpdateMajorTerm(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.MajorTermReqeustUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateMajorTerm(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}
