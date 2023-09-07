package ioc

import (
	"github.com/redis/go-redis/v9"

	"test/webook/internal/repository/cache"
)

func InitCodeCache(client redis.Cmdable) cache.CodeCache {
	return cache.NewRedisCodeCache(client)
}

//func InitCodeCache(c *bigcache.BigCache) cache.CodeCache {
//	return cache.NewBigCodeCache(c)
//}
