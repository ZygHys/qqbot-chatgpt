package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qqbot/api_key"
	"qqbot/service"
)

func Chat4(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.ChatSend(ctx.Query("msg"), service.GPT4, api_key.ApiKey4))
}
