package db

import (
	"context"

	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/busyfree/shorturl-go/util/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoDBs = make(map[string]*mongo.Database)

func GetMongoDB(ctx context.Context, name string) *mongo.Database {
	lock.RLock()
	db := mongoDBs[name]
	lock.RUnlock()

	if db != nil {
		return db
	}
	// mongodb://localhost:27017
	config := conf.GetStrMapStr("MONGODB_" + name + "")
	if config == nil {
		config = conf.GetStrMapStr("MONGODB_DEFAULT")
		if config == nil {
			panic("missing MONGODB config")
		}
	}
	if len(config["dsn"]) == 0 {
		panic("missing MONGODB dsn")
	}
	if len(config["db"]) == 0 {
		panic("missing MONGODB db")
	}
	opts := options.Client().ApplyURI(config["dsn"])
	if len(config["username"]) > 0 && len(config["password"]) > 0 {
		mechanism := "SCRAM-SHA-1"
		if val, ok := config["mechanism"]; ok {
			if len(val) > 0 {
				mechanism = val
			}
		}
		credential := options.Credential{
			AuthMechanism: mechanism,
			Username:      config["username"],
			Password:      config["password"],
		}
		opts = opts.SetAuth(credential)
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Get(ctx).Panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Get(ctx).Panic(err)
	}
	db = client.Database(config["db"])
	lock.Lock()
	mongoDBs[name] = db
	lock.Unlock()
	return db
}

func ResetMongoDB() {
	if !conf.GetBool("HOT_LOAD_DB") {
		return
	}
	lock.Lock()
	oldDBs := mongoDBs
	dbs = make(map[string]*DB, 4)
	lock.Unlock()

	for k, db := range oldDBs {
		_ = db.Client().Disconnect(context.TODO())
		delete(dbs, k)
	}
	return
}
