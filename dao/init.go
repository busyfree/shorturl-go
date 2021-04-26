package dao

import (
	"context"
	"strings"

	"github.com/busyfree/shorturl-go/models/schema"
	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/busyfree/shorturl-go/util/db"
)

var (
	tableLink = new(schema.Link)
)

func SyncXORMTables() {
	if strings.ToLower(conf.GetString("DATA_STORAGE")) == "mysql" {
		ctx := context.Background()
		c := db.GetXORM(ctx, "default")
		_ = c.Sync2(tableLink)
	}
}
