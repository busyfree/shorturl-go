

package webgin

import (
	"context"
	"github.com/busyfree/shorturl-go/util/log"
	"github.com/gin-gonic/gin"
)

const (
	BASEURL = ""
)

var (
	GinRoute *gin.Engine
	logger   = log.Get(context.Background())
)

func InitWebGin() {
	GinRoute = gin.New()
	initRoute()
}
