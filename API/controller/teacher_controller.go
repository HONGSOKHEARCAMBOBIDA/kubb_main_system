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

type TeacherController struct {
	service service.TeacherService
}

func NewTeacherController() TeacherController {
	return TeacherController{
		service: service.NewTeacherService(),
	}
}

func (cr *TeacherController) CreateTeacherRate(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	var input request.TeacherRateRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateTeacherRate(c.Request.Context(), input, userID); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *TeacherController) CreateTeacher(c *gin.Context) {
	var input request.TeacherRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateTeacher(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *TeacherController) GetTeacherFilter(c *gin.Context) {

	filter := map[string]string{
		"faculty_id": c.Query("faculty_id"),
	}

	data, err := cr.service.GetTeacherFilter(
		c.Request.Context(),
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

	for i := range data {
		data[i].Dob = helper.FormatDate(data[i].Dob)
	}

	share.RespondDate(c, http.StatusOK, data)
}

func (cr *TeacherController) GetTeacher(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"name": c.Query("name"),
	}

	data, meta, err := cr.service.GetTeacher(
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

	for i := range data {
		data[i].Dob = helper.FormatDate(data[i].Dob)
	}

	share.ResponsePagination(c, http.StatusOK, data, meta)
}

func (cr *TeacherController) UpdateTeacher(c *gin.Context) {
	uuid, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.TeacherRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateTeacher(c.Request.Context(), uuid, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *TeacherController) Toggle(c *gin.Context) {
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
