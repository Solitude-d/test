package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"test/webook/config"
	"test/webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3009)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
