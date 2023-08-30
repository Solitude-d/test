package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//
	//rdb := initRedis()
	//user := initUser(db, rdb)
	//user.UserRouteRegister(server)
	//server := gin.Default()

	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello mac")
	})
	//server.Run(":8080") // 监听并在 :8080 上启动服务
	server.Run(":8081") // 监听并在 :8081 上启动服务
}

//func initDB() *gorm.DB {
//	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3009)/webook"))
//	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

//func initRedis() redis.Cmdable {
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: config.Config.Redis.Addr,
//	})
//	return redisClient
//}

//func initWebServer() *gin.Engine {
//	server := gin.Default()
//
//	//cmd := redis.NewClient(&redis.Options{
//	//	Addr:     "webook-redis:16379",
//	//	Password: "",
//	//	DB:       1,
//	//})
//
//	//cmd := redis.NewClient(&redis.Options{
//	//	Addr:     config.Config.Redis.Addr,
//	//	Password: "",
//	//	DB:       1,
//	//})
//
//	//一秒钟 100个请求
//	//server.Use(ratelimit.NewBuilder(cmd, time.Second, 100).Build())
//
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins: []string{"localhost:3000"},
//		//AllowMethods:     []string{"PUT", "PATCH"},
//		AllowHeaders:     []string{"authorization", "content-type"},
//		ExposeHeaders:    []string{"x-jwt-token"},
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			//如果origin 包含 http://localhost 则接收请求
//			if strings.HasPrefix(origin, "http://localhost") {
//				return true
//			}
//			return strings.Contains(origin, "webook.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//
//	//为服务添加session
//	//store := cookie.NewStore([]byte("secret"))
//
//	//第一个参数 最大空闲连接数量
//	//store, err := sessRedis.NewStore(16, "tcp", "localhost:6379",
//	//	"", []byte("xHd&^OrleeXM@Yq40gfww%8S%eND1*md"), []byte("O$$f20qm05iP1tcYqT1$pcB15v3L@4Iv"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	//server.Use(sessions.Sessions("webookses", store))
//
//	//server.Use(middleware.NewLoginMiddlewareBuilder().
//	//	IgnorePaths("/users/login", "/users/signup").Builder())
//
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
//		IgnorePaths("/hello", "/users/login", "/users/signup",
//			"users/login_sms/code/send", "users/login_sms").Builder())
//	return server
//}

//func initUser(db *gorm.DB, rdb redis.Cmdable) *user2.UserHandler {
//	udao := dao.NewUserDao(db)
//	ud := cache.NewUserCache(rdb)
//	repo := repository.NewUserRepository(udao, ud)
//	svc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rdb)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	//smsSvc := tencent.NewService()
//	smsSvc := memory.NewService()
//	codeSvc := service.NewCodeService(codeRepo, smsSvc)
//	user := user2.NewUserHandler(svc, codeSvc)
//	return user
//}
