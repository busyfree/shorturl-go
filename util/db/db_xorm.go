package db

import (
	"context"
	"github.com/spf13/cast"
	"time"

	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/busyfree/shorturl-go/util/log"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var dbXORMs = make(map[string]*xorm.Engine, 4)

func GetXORM(ctx context.Context, name string) *xorm.Engine {
	lock.RLock()
	db := dbXORMs[name]
	lock.RUnlock()

	if db != nil {
		return db
	}

	config := conf.GetStrMapStr("DB_" + name)
	sqldb, err := xorm.NewEngine("mysql", config["dsn"])
	if err != nil {
		log.Get(ctx).Panic(err)
	}
	// 不能设太多，数据库最大连接数总共不宜超过 2k
	maxCon := cast.ToInt(config["max_open_conns"])
	maxIdle := cast.ToInt(config["max_idle_conns"])
	maxLife := cast.ToDuration(config["max_life"])

	if maxCon <= 0 {
		maxCon = 10
	}
	if maxIdle <= 0 {
		maxIdle = 10
	}
	if maxLife <= 0 {
		maxLife = 60
	}
	sqldb.SetMaxOpenConns(maxCon)
	sqldb.SetMaxIdleConns(maxIdle)
	sqldb.SetConnMaxLifetime(maxLife * time.Second)
	_ = sqldb.Ping()
	lock.Lock()
	dbXORMs[name] = sqldb
	lock.Unlock()
	return sqldb
}

// ResetXORM 关闭所有 DB 连接
// 新调用 GetXORM 方法时会使用最新 DB 配置创建连接
func ResetXORM() {
	for k, db := range dbXORMs {
		db.Close()
		delete(dbXORMs, k)
	}
}
