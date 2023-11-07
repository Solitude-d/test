package repository

import (
	"context"

	"github.com/gin-gonic/gin"

	"test/webook/internal/domain"
	"test/webook/internal/repository/cache"
	"test/webook/internal/repository/dao"
	"test/webook/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLike(ctx *gin.Context, biz string, bizId int64, uid int64) error
	IncrCancelLike(ctx *gin.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx *gin.Context, biz string, bizId int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, bizId int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bizId int64, uid int64) (bool, error)
}

type CachedReadCntRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.Logger
}

func (c *CachedReadCntRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}
	daoIntr, err := c.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = c.toDomain(daoIntr)
	go func() {
		err := c.cache.Set(ctx, biz, bizId, intr)
		if err != nil {
			c.l.Error("回写缓存失败",
				logger.Strings("biz", biz),
				logger.Int64("bizId", bizId),
			)
		}
	}()
	return intr, nil
}

func (c *CachedReadCntRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}

func (c *CachedReadCntRepository) Liked(ctx context.Context, biz string, bizId int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, err
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) Collected(ctx context.Context, biz string, bizId int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectionInfo(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, err
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) AddCollectionItem(ctx *gin.Context, biz string, bizId int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Cid:   cid,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) IncrCancelLike(ctx *gin.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) IncrLike(ctx *gin.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, id)
	if err != nil {
		return err
	}
	return c.cache.IncrReadCntIfPresent(ctx, biz, id)
}

func NewInteractiveRepository(d dao.InteractiveDAO,
	cache cache.InteractiveCache,
	l logger.Logger) InteractiveRepository {
	return &CachedReadCntRepository{
		dao:   d,
		cache: cache,
		l:     l,
	}
}
