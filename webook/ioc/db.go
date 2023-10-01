package ioc

import (
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"test/webook/internal/repository/dao"
	"test/webook/pkg/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
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
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		//这里缺少一个writer
		Logger: glogger.New(gprmLoggerFun(l.Debug), glogger.Config{
			//慢查询阈值 只有执行时间超过这个阈值 才会使用logger
			SlowThreshold: time.Millisecond * 10,
			//true 会原生sql
			ParameterizedQueries: false,
			//是否忽略 数据库里没数据
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gprmLoggerFun func(msg string, fields ...logger.Field)

func (g gprmLoggerFun) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: msg, Value: args})
}
