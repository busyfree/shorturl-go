package acp

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/busyfree/shorturl-go/service"
	"github.com/busyfree/shorturl-go/util/crypto/md5"
	"github.com/busyfree/shorturl-go/util/shorturl"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkController struct{}

type ShortUrlForm struct {
	Link    string `form:"url" json:"url" xml:"url"  binding:"required"`
	Project string `form:"p" json:"p" xml:"p" binding:"-"`
}

type QueryForm struct {
	Project string `form:"p" json:"p" xml:"p" binding:"-"`
}

func (c *LinkController) Find(ctx *gin.Context) {
	var form QueryForm
	code := ctx.Param("code")
	if len(code) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": "missing code"})
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
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 1, "message": "ok", "data": linkService.Dao})
	return
}

func (c *LinkController) Save(ctx *gin.Context) {
	var form ShortUrlForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
		return
	}
	link, _ := url.QueryUnescape(form.Link)
	form.Link = link
	uri, err := url.Parse(form.Link)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
		return
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": "url scheme error, should be http or https."})
		return
	}
	linkService := service.NewLinkService()
	linkService.Dao.Url = form.Link
	if len(form.Project) == 0 {
		linkService.Dao.Project = "default"
	}
	err = linkService.FindOneLink(ctx)
	if err != nil {
		if err != sql.ErrNoRows && err != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
			return
		}
	}
	if !linkService.Dao.ID_.IsZero() || linkService.Dao.Id > 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": 1, "message": "ok", "data": linkService.Dao})
		return
	}
	linkService.Dao.IP = ctx.ClientIP()
	linkService.Dao.Hash = md5.EncryptString(form.Link)
	linkService.Dao.Key = shorturl.ShortUrl(form.Link)
	err = linkService.SaveOneLink(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 1, "message": "OK", "data": linkService.Dao})
	return
}
