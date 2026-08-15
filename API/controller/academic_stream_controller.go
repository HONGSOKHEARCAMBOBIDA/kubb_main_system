package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AcademicStreamController struct {
	service service.AcademicStreamService
}

func NewAcademicStreamController() AcademicStreamController {
	return AcademicStreamController{
		service: service.NewAcademicStreamService(),
	}
}

func (cr *AcademicStreamController) GetAcademicStream(c *gin.Context) {
	data, err := cr.service.GetAcademicStream(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
