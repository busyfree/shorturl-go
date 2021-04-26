package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/busyfree/shorturl-go/models/schema"
	"github.com/busyfree/shorturl-go/util/crypto/md5"
	"github.com/busyfree/shorturl-go/util/ctxkit"
	"github.com/busyfree/shorturl-go/util/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkDao struct {
	schema.Link
}

func NewLinkDao() *LinkDao {
	return new(LinkDao)
}

func (dao *LinkDao) FindOneLinkByHash(ctx context.Context) (err error) {
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	hash := md5.EncryptString(dao.Url)
	sqlStr := "SELECT id, code, url, hash, project, ip, created_at FROM %s"
	whereStr := " WHERE hash=? AND deleted_at=0"
	whereValues := make([]interface{}, 0, 0)
	whereValues = append(whereValues, hash)
	if len(dao.Project) > 0 {
		whereStr += " AND project=?"
		whereValues = append(whereValues, dao.Project)
	}
	sqlSelect := fmt.Sprintf(sqlStr+whereStr, dao.TableName(ctx))
	q := db.SQLSelect(dao.TableName(ctx), sqlSelect)
	result := c.QueryRowContext(ctx, q, whereValues...)
	if result == nil {
		err = sql.ErrNoRows
		return
	}
	err = result.Scan(&dao.Id, &dao.Key, &dao.Url, &dao.Hash, &dao.Project, &dao.IP, &dao.CreatedAt)
	return
}

func (dao *LinkDao) FindOneLinkByCode(ctx context.Context) (err error) {
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlStr := "SELECT id, code, url, hash, project, ip, created_at FROM %s"
	whereStr := " WHERE code=? AND deleted_at=0"
	whereValues := make([]interface{}, 0, 0)
	whereValues = append(whereValues, dao.Key)
	if len(dao.Project) > 0 {
		whereStr += " AND project=?"
		whereValues = append(whereValues, dao.Project)
	}
	sqlSelect := fmt.Sprintf(sqlStr+whereStr, dao.TableName(ctx))
	q := db.SQLSelect(dao.TableName(ctx), sqlSelect)
	result := c.QueryRowContext(ctx, q, whereValues...)
	if result == nil {
		err = sql.ErrNoRows
		return
	}
	err = result.Scan(&dao.Id, &dao.Key, &dao.Url, &dao.Hash, &dao.Project, &dao.IP, &dao.CreatedAt)
	return
}

func (dao *LinkDao) FindOneLinkByHashFromMongoDB(ctx context.Context) (err error) {
	c := db.GetMongoDB(ctx, ctxkit.GetProjectDBName(ctx))
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	md5Str := md5.EncryptString(dao.Url)
	whereValues := bson.D{{"h", md5Str}}
	if len(dao.Project) > 0 {
		whereValues = append(whereValues, bson.E{Key: "p", Value: dao.Project})
	}
	row := c.Collection(dao.CollectionName(ctx)).FindOne(timeoutCtx, whereValues)
	if row.Err() != nil {
		if row.Err() != mongo.ErrNoDocuments {
			err = row.Err()
			return
		}
	}
	var rs bson.D
	err = row.Decode(&rs)
	if err != nil {
		return
	}
	err = dao.ToStruct(rs)
	if err != nil {
		return
	}
	return
}

func (dao *LinkDao) FindOneLinkByCodeFromMongoDB(ctx context.Context) (err error) {
	c := db.GetMongoDB(ctx, ctxkit.GetProjectDBName(ctx))
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	whereValues := bson.D{{"k", dao.Key}}
	if len(dao.Project) > 0 {
		whereValues = append(whereValues, bson.E{Key: "p", Value: dao.Project})
	}
	row := c.Collection(dao.CollectionName(ctx)).FindOne(timeoutCtx, whereValues)
	if row.Err() != nil {
		if row.Err() != mongo.ErrNoDocuments {
			err = row.Err()
			return
		}
		err = nil
	}
	var rs bson.D
	err = row.Decode(&rs)
	if err != nil {
		return
	}
	err = dao.ToStruct(rs)
	if err != nil {
		return
	}
	return
}

func (dao *LinkDao) SaveOneLinkToMongoDB(ctx context.Context) (err error) {
	c := db.GetMongoDB(ctx, ctxkit.GetProjectDBName(ctx))
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	id, err := c.Collection(dao.CollectionName(ctx)).InsertOne(timeoutCtx, dao.ToMongoDocument())
	if err != nil {
		return
	}
	if i, ok := id.InsertedID.(primitive.ObjectID); ok {
		dao.ID_ = i
	}
	return
}

func (dao *LinkDao) SaveOneLink(ctx context.Context) (err error) {
	dao.BeforeInsert()
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlInsert := fmt.Sprintf("INSERT INTO %s(id, code, url, hash, project, ip, created_at, deleted_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", dao.TableName(ctx))
	q := db.SQLInsert(dao.TableName(ctx), sqlInsert)
	_, err = c.ExecContext(
		ctx,
		q,
		dao.Id,
		dao.Key,
		dao.Url,
		dao.Hash,
		dao.Project,
		dao.IP,
		dao.CreatedAt,
		0,
		dao.UpdatedAt)
	dao.ConvertIdsToStr()
	return
}
