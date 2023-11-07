package startup

import (
	"test/webook/internal/service/oauth2/wechat"
	"test/webook/pkg/logger"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.Logger) wechat.Service {
	return wechat.NewService("", "")
}
