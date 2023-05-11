package filter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var limit = 0
var lastUpdate = time.Now()
var mu sync.Mutex

//todo 升级桶限流
func CurrentFilter(ctx *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	// 简陋限流器
	if time.Since(lastUpdate) > 60*time.Second {
		limit = 0
		lastUpdate = time.Now()
	}
	if limit > 30 {
		ctx.JSON(http.StatusOK, "等一会......限流了")
		ctx.Abort()
	}
	limit++
	ctx.Next()
}
