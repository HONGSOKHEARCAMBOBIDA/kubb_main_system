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

type ClassCurriculumnController struct {
	service service.ClassCurriculumnService
}

func NewClassCurriculumnController() ClassCurriculumnController {
	return ClassCurriculumnController{
		service: service.NewClassCurriculumnService(),
	}
}

func (cr *ClassCurriculumnController) CreateClassCurriculumn(c *gin.Context) {
	var input request.ClassCurriculumnRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateClassCurriculumn(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *ClassCurriculumnController) GetClassCurriculumn(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	filter := map[string]string{
		"programme_id":  c.Query("programme_id"),
		"faculty_id":    c.Query("faculty_id"),
		"department_id": c.Query("department_id"),
		"term_id":       c.Query("term_id"),
		"generation_id": c.Query("generation_id"),
		"academic_id":   c.Query("academic_id"),
	}
	data, meta, err := cr.service.GetClassCurriculumn(
		c.Request.Context(),
		request.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		filter,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}

		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponsePagination(c, http.StatusOK, data, meta)
}

func (cr *ClassCurriculumnController) GetClassCurriculumnWithTeacherRate(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	filter := map[string]string{
		"programme_id":  c.Query("programme_id"),
		"faculty_id":    c.Query("faculty_id"),
		"department_id": c.Query("department_id"),
		"term_id":       c.Query("term_id"),
		"generation_id": c.Query("generation_id"),
		"academic_id":   c.Query("academic_id"),
	}
	data, meta, err := cr.service.GetClassCurriculumnWithTeacherRate(
		c.Request.Context(),
		request.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		filter,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}

		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponsePagination(c, http.StatusOK, data, meta)
}

func (cr *ClassCurriculumnController) UpdateClassCurriculumn(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.ClassCurriculumnRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	//log.Printf("Update request: %+v", input)
	if err := cr.service.UpdateClassCurriculumn(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}
