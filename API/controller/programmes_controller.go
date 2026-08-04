package controller

import (
	"mysql/constant/share"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Programmescontroller struct {
	service service.ProgramesService
}

func NewProgrammescontroller() Programmescontroller {
	return Programmescontroller{
		service: service.NewProgrammesService(),
	}
}

func (cr *Programmescontroller) GetProgrammes(c *gin.Context) {
	data, err := cr.service.GetProgrammes(c)
	if err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
