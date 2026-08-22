package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
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
