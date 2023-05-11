package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qqbot/api_key"
	"qqbot/service"
)

func Chat(c *gin.Context) {
	c.JSON(http.StatusOK, service.ChatSend(c.Query("msg"), service.GPT3_5, api_key.ApiKey3))
}
