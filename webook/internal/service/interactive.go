package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"test/webook/internal/domain"
	"test/webook/internal/repository"
	"test/webook/pkg/logger"
)

type Interactive interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx *gin.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx *gin.Context, biz string, bizId int64, uid int64) error
	Collect(ctx *gin.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
	l    logger.Logger
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	var (
		eg        errgroup.Group
		intr      domain.Interactive
		liked     bool
		collected bool
	)
	eg.Go(func() error {
		var err error
		intr, err = svc.repo.Get(ctx, biz, bizId)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = svc.repo.Liked(ctx, biz, bizId, uid)
		return err
	})

	eg.Go(func() error {
		var err error
		collected, err = svc.repo.Collected(ctx, biz, bizId, uid)
		return err
	})
	err := eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Collected = collected
	return intr, nil
}

func (svc *interactiveService) Collect(ctx *gin.Context, biz string, bizId, cid, uid int64) error {
	return svc.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func (svc *interactiveService) Like(ctx *gin.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.IncrLike(ctx, biz, bizId, uid)
}

func (svc *interactiveService) CancelLike(ctx *gin.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.IncrCancelLike(ctx, biz, bizId, uid)
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func NewInteractiveService(repo repository.InteractiveRepository,
	l logger.Logger) Interactive {
	return &interactiveService{
		repo: repo,
		l:    l,
	}
}
