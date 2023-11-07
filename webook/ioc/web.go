package ioc

import (
	"context"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	"test/webook/internal/web"
	ijwt "test/webook/internal/web/jwt"
	"test/webook/internal/web/middleware"
	"test/webook/pkg/ginx/middlewares/logger"
	"test/webook/pkg/ginx/middlewares/ratelimit"
	logger2 "test/webook/pkg/logger"
)

func InitWebServer(mdl []gin.HandlerFunc, uhdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler,
	articleHandler *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl...)
	uhdl.UserRouteRegister(server)
	oauth2WechatHdl.RegisterRoutes(server)
	articleHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler,
	l logger2.Logger) []gin.HandlerFunc {
	bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
		//什么环境初始化哪种等级的日志
		l.Debug("Http请求", logger2.Field{Key: "al", Value: al})
	}).AllowReqBody(true).AllowResBody(true)
	//动态监听是否写请求体和响应体日志
	viper.OnConfigChange(func(in fsnotify.Event) {
		req := viper.GetBool("web.reqbody")
		bd.AllowReqBody(req)

		res := viper.GetBool("web.resbody")
		bd.AllowResBody(res)

	})
	return []gin.HandlerFunc{
		corsHel(),
		bd.Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/hello", "/users/login", "/users/signup",
				"users/login_sms/code/send", "users/login_sms", "oauth2/wechat/authurl",
				"oauth2/wechat/callback", "users/refresh_token").Builder(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHel() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"localhost:3000"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"authorization", "content-type"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//如果origin 包含 http://localhost 则接收请求
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "webook.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
