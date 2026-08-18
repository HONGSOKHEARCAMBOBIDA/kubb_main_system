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

type SchoolarshipController struct {
	service service.SchoolarshipService
}

func NewSchoolarshipController() SchoolarshipController {
	return SchoolarshipController{
		service: service.NewSchoolarshipService(),
	}
}

func (cr *SchoolarshipController) GetSchoolarship(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)

	filter := map[string]string{
		"name": c.Query("name"),
		"code": c.Query("code"),
	}

	data, meta, err := cr.service.GetSchoolarship(
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

			share.ResponseError(
				c,
				http.StatusGatewayTimeout,
				err.Error(),
			)
			return
		}

		share.ResponseError(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	share.ResponsePagination(c, http.StatusOK, data, meta)
}

func (cr *SchoolarshipController) CreateSchoolarship(c *gin.Context) {
	var input request.SchoolarshipRequestCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := cr.service.CreateSchoolarship(
		c.Request.Context(),
		input,
	); err != nil {
		share.RespondServiceError(c, err)
		return
	}

	share.ResponseSuccess(
		c,
		http.StatusOK,
		share.Created,
	)
}

func (cr *SchoolarshipController) UpdateSchoolarship(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		return
	}

	var input request.SchoolarshipRequestUpdate

	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := cr.service.UpdateSchoolarship(
		c.Request.Context(),
		id,
		input,
	); err != nil {
		share.RespondServiceError(c, err)
		return
	}

	share.ResponseSuccess(
		c,
		http.StatusOK,
		share.Updated,
	)
}

func (cr *SchoolarshipController) Toggle(c *gin.Context) {
	id, ok := utils.GetParamUUID(c)
	if !ok {
		share.ResponseError(
			c,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}

	if err := cr.service.Toggle(
		c.Request.Context(),
		id,
	); err != nil {
		share.RespondServiceError(c, err)
		return
	}

	share.ResponseSuccess(
		c,
		http.StatusOK,
		share.Updated,
	)
}
