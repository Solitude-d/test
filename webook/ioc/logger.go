package ioc

import (
	"go.uber.org/zap"

	"test/webook/pkg/logger"
)

//func InitLogger() *zap.Logger {
//	l, err := zap.NewDevelopment()
//	if err != nil {
//		panic(err)
//	}
//	return l
//}

func InitLogger() logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
