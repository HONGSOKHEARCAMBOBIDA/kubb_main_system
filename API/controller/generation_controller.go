package controller

import (
	"mysql/constant/share"
	"mysql/service"
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
