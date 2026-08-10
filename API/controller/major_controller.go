package controller

import (
	"context"
	"errors"
	"log"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MajorController struct {
	service service.MajorService
}

func NewMajorController() MajorController {
	return MajorController{
		service: service.NewMajorService(),
	}
}

func (cr *MajorController) GetMajor(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"programme_id":  c.Query("programme_id"),
		"faculty_id":    c.Query("faculty_id"),
		"department_id": c.Query("department_id"),
	}

	data, meta, err := cr.service.GetMajor(c.Request.Context(), request.Pagination{
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

func (cr *MajorController) GetMajorByDepartment(c *gin.Context) {
	departmentID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetMajorByDepartment(c, departmentID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *MajorController) CreateMajor(c *gin.Context) {
	var input request.MajorRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf(err.Error())
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.CreateMajor(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *MajorController) UpdateMajor(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	var input request.MajorRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.UpdateMajor(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *MajorController) Toggle(c *gin.Context) {
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
