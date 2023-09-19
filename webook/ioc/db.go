package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"test/webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3009)/webook"))
	//dsn := viper.GetString("db.mysql.dsn")
	//db, err := gorm.Open(mysql.Open(dsn))
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	//err := viper.UnmarshalKey("db.mysql", &cfg)
	// remote 不支持 db.mysql
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
