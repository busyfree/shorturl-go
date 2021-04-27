package dao

import (
	"context"
	"strings"

	"github.com/busyfree/shorturl-go/models/schema"
	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/busyfree/shorturl-go/util/crypto/md5"
	"github.com/busyfree/shorturl-go/util/db"
	"github.com/busyfree/shorturl-go/util/shorturl"
)

var (
	tableLink = new(schema.Link)
)

func SyncXORMTables() {
	ctx := context.Background()
	if strings.ToLower(conf.GetString("DATA_STORAGE")) == "mysql" {
		c := db.GetXORM(ctx, "default")
		_ = c.Sync2(tableLink)
	} else {
		linkDao := NewLinkDao()
		linkDao.Url = "https://golang.org"
		linkDao.Hash = md5.EncryptString(linkDao.Url)
		linkDao.Key = shorturl.ShortUrl(linkDao.Url)
		linkDao.IP = "127.0.0.1"
		linkDao.Project = "default"
		_ = linkDao.SaveOneLinkToMongoDB(ctx)
		_, _ = linkDao.CreateMongoDBIndex(ctx)
	}
}
