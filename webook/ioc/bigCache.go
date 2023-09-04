package ioc

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

func InitBigCache() *bigcache.BigCache {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Minute*10))
	return cache
}
