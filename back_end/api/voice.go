package api

import (
	"net/http"

	"github.com/Chabuduo04/mate/back_end/services"
	"github.com/gin-gonic/gin"
)

type VoiceApi struct {
	Service services.TTSService
}

func NewVoiceApi(service services.TTSService) VoiceApi {
	return VoiceApi{Service: service}
}

func (v *VoiceApi) VoiceList(c *gin.Context) {
	body, err := v.Service.ListVoicesRaw()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list voices error"})
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
