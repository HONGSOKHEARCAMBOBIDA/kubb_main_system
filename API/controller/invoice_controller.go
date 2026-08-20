package controller

import (
	"mysql/constant/share"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvoiceController struct {
	service service.InvoiceService
}

func NewInvoiceController() InvoiceController {
	return InvoiceController{
		service: service.NewInvoiceService(),
	}
}

func (cr *InvoiceController) CreateInvoice(c *gin.Context) {
	var input request.InvoiceRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateInvoice(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}
