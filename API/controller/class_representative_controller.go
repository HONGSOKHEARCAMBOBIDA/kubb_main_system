package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassRepresentativeController struct {
	service service.ClassRepresentativeService
}

func NewClassRepresentativeController() ClassRepresentativeController {
	return ClassRepresentativeController{
		service: service.NewClassRepresentative(),
	}
}

func (cr *ClassRepresentativeController) CreateClassRepresentative(c *gin.Context) {
	var input request.ClassRepresentativeRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateClassRepresentative(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}
