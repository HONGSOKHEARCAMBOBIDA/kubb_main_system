package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassOfferingController struct {
	service service.ClassOfferingService
}

func NewClassOfferingController() ClassOfferingController {
	return ClassOfferingController{
		service: service.NewClassOfferingService(),
	}
}

func (cr *ClassOfferingController) CreateClassOffering(c *gin.Context) {
	var input request.ClassOfferingRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateClassOffering(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}
