package api

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/busyfree/shorturl-go/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkController struct{}

type QueryForm struct {
	Project string `form:"p" json:"p" xml:"p" binding:"-"`
}

func (c *LinkController) Redirect(ctx *gin.Context) {
	var form QueryForm
	code := ctx.Param("code")
	if len(code) == 0 {
		ctx.String(http.StatusOK, "welcome to short url service")
		return
	}
	if err := ctx.ShouldBindQuery(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
		return
	}
	linkService := service.NewLinkService()
	linkService.Dao.Key = code
	if len(form.Project) == 0 {
		linkService.Dao.Project = "default"
	}
	err := linkService.FindOneLinkByKey(ctx)
	if err != nil {
		if err == sql.ErrNoRows || err == mongo.ErrNoDocuments {
			ctx.String(http.StatusOK, "welcome to short url service")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
		return
	}
	if len(linkService.Dao.Url) == 0 {
		ctx.String(http.StatusOK, "welcome to short url service")
		return
	}
	target, errUnescape := url.QueryUnescape(linkService.Dao.Url)
	if errUnescape != nil {
		ctx.String(http.StatusOK, errUnescape.Error())
		return
	}
	//@todo write some analysis data to nsq or other async queue
	ctx.Redirect(http.StatusFound, target)
	return
}