package ioc

import (
	"github.com/allegro/bigcache/v3"

	"test/webook/internal/repository/cache"
)

//func InitCodeCache(client redis.Cmdable) cache.CodeCache {
//	return cache.NewRedisCodeCache(client)
//}

func InitCodeCache(c *bigcache.BigCache) cache.CodeCache {
	return cache.NewBigCodeCache(c)
}
