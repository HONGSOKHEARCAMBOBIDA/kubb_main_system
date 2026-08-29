package controller

import (
	"context"
	"errors"
	"log"

	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
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

func (cr *StudentController) GetStudentCategory(c *gin.Context) {
	data, err := cr.service.GetStudentCategory(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *StudentController) CreateStudent(c *gin.Context) {
	var input request.StudentRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("CreateStudent request: %+v", err)
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateStudent(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *StudentController) GetCourseRegistration(c *gin.Context) {
	filter := map[string]string{
		"class_offering_id": c.Query("class_offering_id"),
		"attendance_id":     c.Query("attendance_id"),
	}
	data, err := cr.service.GetCourseRegistration(c, filter)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}

		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *StudentController) GetStudent(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"name":                c.Query("name"),
		"student_category_id": c.Query("student_category_id"),
		"group_id":            c.Query("group_id"),
		"phone":               c.Query("phone"),
		"stream_id":           c.Query("stream_id"),
	}

	data, meta, err := cr.service.GetStudent(
		c.Request.Context(),
		request.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		filter,
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}

		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	for i := range data {
		data[i].DateOfBirth = helper.FormatDate(data[i].DateOfBirth)
		for j := range data[i].StudentEducation {
			data[i].StudentEducation[j].StartDate =
				helper.FormatDate(data[i].StudentEducation[j].StartDate)

			data[i].StudentEducation[j].EndDate =
				helper.FormatDate(data[i].StudentEducation[j].EndDate)

			data[i].StudentEducation[j].CertificateDate =
				helper.FormatDate(data[i].StudentEducation[j].CertificateDate)
		}
	}

	share.ResponsePagination(c, http.StatusOK, data, meta)
}

func (cr *StudentController) UpdateStudent(c *gin.Context) {
	id, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	var input request.StudentRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	//log.Printf("Update request: %+v", input)
	if err := cr.service.UpdateStudent(c.Request.Context(), id, input, userID); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}
