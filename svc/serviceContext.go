package svc

import (
	"SQLsync/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ServiceContext struct {
	Config *config.Config
	DB     *sqlx.DB
}

func NewServiceContext(conf *config.Config) *ServiceContext {
	return &ServiceContext{
		Config: conf,
		DB:     sqlx.MustConnect(conf.DataSource.Type, conf.DataSource.Url),
	}
}
