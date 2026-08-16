package controller

import (
	"context"
	"errors"
	"log"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	service service.StudentService
}

func NewStudentController() StudentController {
	return StudentController{
		service: service.NewStudentService(),
	}
}

func (cr *StudentController) CreateStudent(c *gin.Context) {
	var input request.StudentRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("CreateStudent request: %+v", input)
	if err := cr.service.CreateStudent(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *StudentController) GetStudent(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"name": c.Query("name"),
	}

	data, meta, err := cr.service.GetStudent(c.Request.Context(), request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)
	for i := range data {
		data[i].DateOfBirth = helper.FormatDate(data[i].DateOfBirth)
		data[i].StudentEducation[i].StartDate = helper.FormatDate(data[i].StudentEducation[i].StartDate)
		data[i].StudentEducation[i].EndDate = helper.FormatDate(data[i].StudentEducation[i].EndDate)
		data[i].StudentEducation[i].CertificateDate = helper.FormatDate(data[i].StudentEducation[i].CertificateDate)
	}

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
