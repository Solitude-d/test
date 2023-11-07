package startup

import (
	"test/webook/pkg/logger"
)

func InitLog() logger.Logger {
	return &logger.NopLogger{}
}
