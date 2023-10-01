//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"test/webook/internal/repository"
	"test/webook/internal/repository/cache"
	"test/webook/internal/repository/dao"
	"test/webook/internal/service"
	"test/webook/internal/web"
	ijwt "test/webook/internal/web/jwt"
	"test/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(

		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		dao.NewUserDao,

		cache.NewUserCache,
		ioc.InitCodeCache,

		//ioc.InitBigCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ioc.NewWeChatHandlerConfig,
		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return gin.Default()
}
