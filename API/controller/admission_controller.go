package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdmissionController struct {
	service service.AdmissionService
}

func NewAdmissionController() AdmissionController {
	return AdmissionController{
		service: service.NewAdmissionService(),
	}
}

func (cr *AdmissionController) GetStudentTermFilter(c *gin.Context) {
	filter := map[string]string{
		"semester_id":   c.Query("semester_id"),
		"study_year_id": c.Query("study_year_id"),
		"term_id":       c.Query("term_id"),
		"major_id":      c.Query("major_id"),
	}
	data, err := cr.service.GetStudentTermFilter(c.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, 200, data)
}

func (cr *AdmissionController) GetAdmission(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"student_id":    c.Query("student_id"),
		"student_name":  c.Query("student_name"),
		"academic_id":   c.Query("academic_id"),
		"generation_id": c.Query("generation_id"),
		"term_id":       c.Query("term_id"),
	}

	data, meta, err := cr.service.GetAdmission(c.Request.Context(), request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)
	for i := range data {
		data[i].Date = helper.FormatDate(data[i].Date)
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

func (cr *AdmissionController) CreateStudentTerm(c *gin.Context) {
	var input request.StudentTermRequestv2
	if err := c.ShouldBindJSON(&input); err != nil {
		// log.Printf("CreateStudent request: %+v", err)
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateStudentTerm(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *AdmissionController) CreateEnrollment(c *gin.Context) {
	var input request.EnrollmentRequestCreateV2
	if err := c.ShouldBindJSON(&input); err != nil {
		// log.Printf("CreateStudent request: %+v", err)
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cr.service.CreateEnrollment(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *AdmissionController) UpdateAdmission(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.AdmissionRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	//log.Printf("Update request: %+v", input)
	if err := cr.service.UpdateAdmission(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *AdmissionController) UpdateEnrollment(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.EnrollmentRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	//log.Printf("Update request: %+v", input)
	if err := cr.service.UpdateEnrollment(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *AdmissionController) UpdateStudentTerm(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.StudentTermRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	//log.Printf("Update request: %+v", input)
	if err := cr.service.UpdateStudentTerm(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *AdmissionController) DeleteEnrollment(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}
	if err := cr.service.DeleteEnrollment(c.Request.Context(), id); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}
