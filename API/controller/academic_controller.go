package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AcademicController struct {
	service service.AcademicService
}

func NewAcademicController() AcademicController {
	return AcademicController{
		service: service.NewAcademicService(),
	}
}

func (cr *AcademicController) GetAcademic(c *gin.Context) {
	data, err := cr.service.GetAcademic(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AcademicController) CreateAcademic(c *gin.Context) {
	var input request.AcademicRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if err := cr.service.CreateAcademic(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}

	share.RespondDate(c, http.StatusCreated, share.Created)
}

func (cr *AcademicController) UpdateAcademic(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.AcademicRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := cr.service.UpdateAcademic(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}

	share.RespondDate(c, http.StatusOK, share.Updated)
}

func (cr *AcademicController) Toggle(c *gin.Context) {
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
