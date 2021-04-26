package service

import (
	"context"
	"strings"

	"github.com/busyfree/shorturl-go/dao"
	"github.com/busyfree/shorturl-go/util/conf"
)

type LinkService struct {
	Dao *dao.LinkDao
}

func NewLinkService() *LinkService {
	s := new(LinkService)
	s.Dao = dao.NewLinkDao()
	return s
}

func (s *LinkService) FindOneLink(ctx context.Context) (err error) {
	if strings.ToLower(conf.GetString("DATA_STORAGE")) == "mysql" {
		return s.FindOneLinkByHash(ctx)
	}
	return s.FindOneLinkByHashFromMongoDB(ctx)
}

func (s *LinkService) FindOneLinkByKey(ctx context.Context) (err error) {
	if strings.ToLower(conf.GetString("DATA_STORAGE")) == "mysql" {
		return s.FindOneLinkByCode(ctx)
	}
	return s.FindOneLinkByCodeFromMongoDB(ctx)
}

func (s *LinkService) SaveOneLink(ctx context.Context) (err error) {
	if strings.ToLower(conf.GetString("DATA_STORAGE")) == "mysql" {
		return s.SaveOneLinkToMySQL(ctx)
	}
	return s.SaveOneLinkToMongoDB(ctx)
}

func (s *LinkService) FindOneLinkByHash(ctx context.Context) (err error) {
	return s.Dao.FindOneLinkByHash(ctx)
}

func (s *LinkService) FindOneLinkByCode(ctx context.Context) (err error) {
	return s.Dao.FindOneLinkByCode(ctx)
}

func (s *LinkService) FindOneLinkByHashFromMongoDB(ctx context.Context) (err error) {
	return s.Dao.FindOneLinkByHashFromMongoDB(ctx)
}

func (s *LinkService) FindOneLinkByCodeFromMongoDB(ctx context.Context) (err error) {
	return s.Dao.FindOneLinkByCodeFromMongoDB(ctx)
}


func (s *LinkService) SaveOneLinkToMongoDB(ctx context.Context) (err error) {
	return s.Dao.SaveOneLinkToMongoDB(ctx)
}

func (s *LinkService) SaveOneLinkToMySQL(ctx context.Context) (err error) {
	return s.Dao.SaveOneLink(ctx)
}
