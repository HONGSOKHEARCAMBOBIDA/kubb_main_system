package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	service service.LocationService
}

func NewLocationController() LocationController {
	return LocationController{
		service: service.NewLocationService(),
	}
}

func (cr *LocationController) GetProvince(c *gin.Context) {
	data, err := cr.service.GetProvince(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *LocationController) GetDistrict(c *gin.Context) {
	departmentID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetDistrict(c, departmentID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *LocationController) GetCommunce(c *gin.Context) {
	departmentID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetCommunce(c, departmentID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *LocationController) GetVillage(c *gin.Context) {
	departmentID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetVillage(c, departmentID)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
