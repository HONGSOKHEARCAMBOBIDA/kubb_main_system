package controller

import (
	"mysql/constant/share"
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
	data, err := cr.service.GetGeneration(c)
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
