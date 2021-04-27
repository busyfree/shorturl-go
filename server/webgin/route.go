package webgin

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/busyfree/shorturl-go/server/webgin/controllers/acp"
	"github.com/busyfree/shorturl-go/server/webgin/controllers/api"
	"github.com/busyfree/shorturl-go/server/webgin/middlewares"
	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/gin-gonic/gin"
)

func formatBool(t bool) string {
	if t {
		return "true"
	}
	return "false"
}

func initRoute() {
	GinRoute.SetFuncMap(template.FuncMap{
		"formatBool": formatBool,
	})

	sessionMaxIdle := conf.GetInt("SESSION_REDIS_MAX_CONS")
	sessionIPPort := conf.GetString("SESSION_REDIS_IP_PORT")
	sessionDB := conf.GetString("SESSION_REDIS_DB")
	if sessionMaxIdle <= 0 {
		sessionMaxIdle = 10
	}
	if len(sessionIPPort) == 0 {
		sessionIPPort = "localhost:6379"
	}
	if len(sessionDB) == 0 {
		sessionDB = "0"
	}

	//GinRoute.LoadHTMLGlob("views/***/**/*")
	fs := filepath.Join(conf.GetConfigPath(), "public", "static")
	GinRoute.StaticFS("/static", http.Dir(fs))
	GinRoute.Use(gin.Recovery())
	middlewares.FilterMiddleware()
	GinRoute.Use(middlewares.Logger())

	GinRoute.MaxMultipartMemory = 100 << 20 //100M

	base := new(api.BaseController)
	link := new(api.LinkController)

	GinRoute.GET("/s/:code", link.Redirect)
	GinRoute.GET(BASEURL+"/ping", base.Ping)

	v1 := GinRoute.Group(BASEURL + "v1")
	v1Dashboard := v1.Group("/acp")
	apiGroup := v1Dashboard.Group("/link")
	{
		apiController := new(acp.LinkController)
		apiGroup.GET("/:code", apiController.Find)
		apiGroup.POST("/save", apiController.Save)
	}
}
