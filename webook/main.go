package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	sessRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"test/webook/internal/repository"
	"test/webook/internal/repository/dao"
	"test/webook/internal/service"
	user2 "test/webook/internal/web"
	"test/webook/internal/web/middleware"
	"test/webook/pkg/ginx/middlewares/ratelimit"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	user := initUser(db)
	user.UserRouteRegister(server)
	server.Run(":8080") // 监听并在 :8080 上启动服务
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	//一秒钟 100个请求
	server.Use(ratelimit.NewBuilder(cmd, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
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
			return strings.Contains(origin, "公司域名")
		},
		MaxAge: 12 * time.Hour,
	}))

	//为服务添加session
	//store := cookie.NewStore([]byte("secret"))

	//第一个参数 最大空闲连接数量
	store, err := sessRedis.NewStore(16, "tcp", "localhost:6379",
		"", []byte("xHd&^OrleeXM@Yq40gfww%8S%eND1*md"), []byte("O$$f20qm05iP1tcYqT1$pcB15v3L@4Iv"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("webookses", store))

	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/login", "/users/signup").Builder())

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().Builder())
	return server
}

func initUser(db *gorm.DB) *user2.UserHandler {
	udao := dao.NewUserDao(db)
	repo := repository.NewUserRepository(udao)
	svc := service.NewUserService(repo)
	user := user2.NewUserHandler(svc)
	return user
}
