package util

import (
	"github.com/busyfree/shorturl-go/util/conf" // init conf
	"github.com/busyfree/shorturl-go/util/db"
	"github.com/busyfree/shorturl-go/util/log"
	"github.com/busyfree/shorturl-go/util/mc"
	"github.com/busyfree/shorturl-go/util/redis"
	"strings"
)

// GatherMetrics 收集一些被动指标
func GatherMetrics() {
	mc.GatherMetrics()
	redis.GatherMetrics()
	db.GatherMetrics()
}

// Reset all utils
func Reset() {
	log.Reset()
	if strings.ToLower(conf.GetString("DATA_STORAGE")) == "mysql" {
		db.ResetXORM()
		db.Reset()
	} else {
		db.ResetMongoDB()
	}
	mc.Reset()
}

// Stop all utils
func Stop() {
}
