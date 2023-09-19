package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
)

type BigCodeCache struct {
	cache *bigcache.BigCache
	lock  sync.Mutex
}

func NewBigCodeCache(cache *bigcache.BigCache) CodeCache {
	return &BigCodeCache{
		cache: cache,
	}
}

func (b *BigCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	key := b.key(biz, phone)
	cntKey := key + ":cnt"
	refreshKey := key + ":fresh"
	val, err := b.cache.Get(key)
	if err != nil {
		return err
	}
	if val == nil {
		//key不存在或者已已经过期
		valMap := map[string]any{
			"code":     code,
			cntKey:     3,
			refreshKey: time.Now().Add(1 * time.Minute),
		}
		byVal, err := json.Marshal(valMap)
		if err != nil {
			return err
		}
		err = b.cache.Set(key, byVal)
		if err != nil {
			return err
		}
	}
	var valMap map[string]any
	err = json.Unmarshal(val, &valMap)
	if err != nil {
		return err
	}
	currentTime := time.Now()
	refTime := valMap[refreshKey]
	if currentTime.After(refTime.(time.Time)) {
		//验证码已经存在，并且已经过了一分钟
		valMap := map[string]any{
			"code":     code,
			cntKey:     3,
			refreshKey: time.Now().Add(1 * time.Minute),
		}
		byVal, err := json.Marshal(valMap)
		if err != nil {
			return err
		}
		err = b.cache.Set(key, byVal)
		if err != nil {
			return err
		}
	} else {
		//验证码存在，但是还没过一分钟
		return ErrCodeSendTooMany
	}
	return nil
}

func (b *BigCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	key := b.key(biz, phone)
	value, err := b.cache.Get(key)
	if err != nil {
		return false, err
	}
	var valMap map[string]any
	err = json.Unmarshal(value, &valMap)
	if err != nil {
		return false, err
	}
	cntKey := key + ":cnt"
	refreshKey := key + ":fresh"
	if valMap[cntKey].(int64) <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}

	if inputCode == valMap["code"].(string) {
		//验证正确
		setMap := map[string]any{
			"code":     inputCode,
			cntKey:     -1,
			refreshKey: time.Now().Add(1 * time.Minute),
		}
		byVal, err := json.Marshal(setMap)
		if err != nil {
			return true, nil
		}
		err = b.cache.Set(key, byVal)
		if err != nil {
			return true, nil
		}
		return true, nil
	} else {
		//验证不正确
		setMap := map[string]any{
			"code":     inputCode,
			cntKey:     valMap[cntKey].(int64) - 1,
			refreshKey: time.Now().Add(1 * time.Minute),
		}
		byVal, err := json.Marshal(setMap)
		if err != nil {
			return false, nil
		}
		err = b.cache.Set(key, byVal)
		if err != nil {
			return false, nil
		}
	}
	return false, ErrUnknowCode
}

func (c *BigCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
