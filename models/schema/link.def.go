package schema

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/busyfree/shorturl-go/util/ctxkit"
	"github.com/busyfree/shorturl-go/util/snowflake"
	"github.com/busyfree/shorturl-go/util/timeutil"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	Id        int64              `xorm:"id BIGINT(20) notnull pk" db:"id" bson:"-" json:"bn_id,omitempty"`
	IdStr     string             `xorm:"-" db:"-" bson:"-" json:"id,omitempty"`
	ID_       primitive.ObjectID `xorm:"-" db:"-" bson:"_id,omitempty" json:"_id,omitempty"`
	Key       string             `xorm:"code CHAR(10) notnull default ''" db:"-" bson:"k,omitempty" json:"k"`
	Url       string             `xorm:"url VARCHAR(255) notnull default ''" db:"-" bson:"u,omitempty" json:"u"`
	Hash      string             `xorm:"hash CHAR(32) notnull default ''" db:"-" bson:"h,omitempty" json:"h"`
	Project   string             `xorm:"project CHAR(20) notnull default ''" db:"-" bson:"p,omitempty" json:"p"`
	IP        string             `xorm:"ip CHAR(49) notnull default ''" db:"-" bson:"ip,omitempty" json:"ip"`
	CreatedAt int64              `xorm:"created_at BIGINT(20) notnull default 0" db:"created_at" bson:"c,omitempty" json:"c,omitempty"`
	DeletedAt int64              `xorm:"deleted_at BIGINT(20) notnull default 0" db:"deleted_at" bson:"-" json:"deleted,omitempty"`
	UpdatedAt int64              `xorm:"updated_at BIGINT(20) notnull default 0" db:"updated_at" bson:"-" json:"updated,omitempty"`
}

func (p *Link) TableName(ctx context.Context) string {
	config := conf.GetStrMapStr("DB_" + ctxkit.GetProjectDBName(ctx))
	prefix := config["prefix"]
	if len(prefix) > 0 {
		return prefix + "_links"
	}
	return "links"
}

func (p *Link) CollectionName(ctx context.Context) string {
	config := conf.GetStrMapStr("MONGODB_" + ctxkit.GetProjectDBName(ctx))
	prefix := config["prefix"]
	if len(prefix) > 0 {
		return prefix + "_links"
	}
	return "links"
}

func (p *Link) ToMongoDocument() bson.D {
	p.BeforeInsert()
	d := bson.D{
		{"_id", primitive.NewObjectID()},
		{"id", p.Id},
		{"k", p.Key},
		{"u", p.Url},
		{"h", p.Hash},
		{"p", p.Project},
		{"ip", p.IP},
		{"c", p.CreatedAt},
	}
	return d
}

func (p *Link) ToStruct(rs bson.D) (err error) {
	if len(rs) == 0 {
		err = errors.New("missing bson.D")
		return
	}
	var jsonMap = make(map[string]interface{}, 0)
	for _, row := range rs {
		if row.Key == "id" {
			jsonMap[row.Key] = cast.ToString(row.Value)
		} else {
			jsonMap[row.Key] = row.Value
		}
	}
	encode, err := json.Marshal(jsonMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(encode, p)
	return
}

func (p *Link) BeforeInsert() {
	if p.Id <= 0 {
		p.Id, _ = snowflake.NextId()
	}
	p.CreatedAt = timeutil.MsTimestampNow()
	p.UpdatedAt = p.CreatedAt
	p.DeletedAt = 0
}

func (p *Link) BeforeUpdate() {
	p.UpdatedAt = timeutil.MsTimestampNow()
}

func (p *Link) ConvertIdsToStr() {
	p.IdStr = cast.ToString(p.Id)
}
