package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"webdevclass-tasks/config"
)

func InitConfigAndDB() (*config.Config, *gorm.DB) {
	cfgFile, err := os.Open("config/config.toml")
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()
	var cfg = config.Read(cfgFile)
	var dbConfig = cfg.DB["Task"]
	var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.UserName, dbConfig.Password, dbConfig.Address, dbConfig.DBName)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return cfg, db
}
