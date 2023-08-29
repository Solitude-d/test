package ioc

import (
	"test/webook/internal/service/sms"
	"test/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
