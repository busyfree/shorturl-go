package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

func (c *BaseController) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
	return
}
