package ioc

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"test/webook/internal/web"
	"test/webook/internal/web/middleware"
	"test/webook/pkg/ginx/middlewares/ratelimit"
)

func InitWebServer(mdl []gin.HandlerFunc, uhdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl...)
	uhdl.UserRouteRegister(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHel(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/hello", "/users/login", "/users/signup",
				"users/login_sms/code/send", "users/login_sms").Builder(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHel() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"localhost:3000"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"authorization", "content-type"},
		ExposeHeaders:    []string{"x-jwt-token"},
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
