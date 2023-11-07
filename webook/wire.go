//go:build wireinject

package main

import (
	"github.com/google/wire"

	"test/webook/internal/events/article"
	"test/webook/internal/repository"
	article2 "test/webook/internal/repository/article"
	"test/webook/internal/repository/cache"
	"test/webook/internal/repository/dao"
	article3 "test/webook/internal/repository/dao/article"
	"test/webook/internal/service"
	"test/webook/internal/web"
	ijwt "test/webook/internal/web/jwt"
	"test/webook/ioc"
)

func InitWebServer() *App {
	wire.Build(

		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewSyncProducer,
		ioc.NewConsumers,

		dao.NewUserDao,
		dao.NewGORMInteractiveDAO,
		article3.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewRedisInteractiveCache,
		ioc.InitCodeCache,

		article.NewInteractiveReadEventConsumer,
		article.NewKafkaProducer,

		//ioc.InitBigCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewInteractiveRepository,
		article2.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ioc.NewWeChatHandlerConfig,
		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
