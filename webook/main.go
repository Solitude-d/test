package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//
	//rdb := initRedis()
	//user := initUser(db, rdb)
	//user.UserRouteRegister(server)
	//server := gin.Default()
	//initViper()
	initViperV1()
	//initViperV2()
	//initViperRemote()
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello mac")
	})
	//server.Run(":8080") // 监听并在 :8080 上启动服务
	server.Run(":8081") // 监听并在 :8081 上启动服务
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	//不replace的话 使用zap.L().  无法打印出日志  直接用前面new出来的logger可以不需要replace
	zap.ReplaceGlobals(logger)

}

func initViper() {
	viper.SetDefault("db.mysql.dsn",
		"root:root@tcp(localhost:13316)/webook")
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	//当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")
	//将配置加载到内存中
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// initViperRemote etcdctl --endpoints=127.0.0.1:12379 put /webook "$(<dev.yaml)"
// etcdctl --endpoints=http://127.0.0.1:12379 get /webook
func initViperRemote() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3",
		"127.0.0.1:12379",
		"/webook")
	if err != nil {
		panic(err)
	}
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}

	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

// initViperReader 一般用于开发调试
func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"
redis:
  addr: "localhost:6379"
`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

// initViperV2 启动时传入配置文件路径 没有则使用设置好的默认值    go run . --config=config/dev.yaml
func initViperV2() {
	cfile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	//默认值
	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:13316)/webook")
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
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
//	//	IgnorePaths("/users/login", "/users/signup").Build())
//
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
//		IgnorePaths("/hello", "/users/login", "/users/signup",
//			"users/login_sms/code/send", "users/login_sms").Build())
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
