package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentTypeController struct {
	service service.DocumentTypeService
}

func NewDocumentTypeController() DocumentTypeController {
	return DocumentTypeController{
		service: service.NewDocumentTypeService(),
	}
}

func (cr *DocumentTypeController) GetDocumentType(c *gin.Context) {
	data, err := cr.service.GetDocumentType(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
