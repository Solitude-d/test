package ioc

import (
	"os"

	"test/webook/internal/service/oauth2/wechat"
	"test/webook/internal/web"
)

func InitWechatService() wechat.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_ID")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	return wechat.NewService(appId, appKey)
}

func NewWeChatHandlerConfig() web.WeChatHandlerConfig {
	return web.WeChatHandlerConfig{
		Secure:   false,
		StateKey: []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixx"),
	}
}
