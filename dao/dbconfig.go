package dao

import (
	"github.com/astaxie/beego"
)

type DBConfig struct {
	Host, User, Pass, DbName string
}

func InitConfig() *DBConfig {

	var DatabaseMySqlHost string
	var DatabaseMySqlUser string
	var DatabaseMySqlPass string
	var DatabaseMySqlDbName string

	if DatabaseMySqlHost = beego.AppConfig.String("DatabaseMySqlHost"); DatabaseMySqlHost == "" {
		DatabaseMySqlHost = "172.18.9.12:13307"
	}

	if DatabaseMySqlUser = beego.AppConfig.String("DatabaseMySqlUser"); DatabaseMySqlUser == "" {
		DatabaseMySqlUser = "cattle"
	}

	if DatabaseMySqlPass = beego.AppConfig.String("DatabaseMySqlPass"); DatabaseMySqlPass == "" {
		DatabaseMySqlPass = "cattle"
	}

	if DatabaseMySqlDbName = beego.AppConfig.String("DatabaseMySqlDbName"); DatabaseMySqlDbName == "" {
		DatabaseMySqlDbName = "cattle"
	}

	return &DBConfig{DatabaseMySqlHost, DatabaseMySqlUser, DatabaseMySqlPass, DatabaseMySqlDbName}
}
