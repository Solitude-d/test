package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"test/webook/internal/domain"
	events "test/webook/internal/events/article"
	"test/webook/internal/repository/article"
	"test/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx *gin.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx *gin.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx *gin.Context, id, uid int64) (domain.Article, error)
}

type artService struct {
	repo article.ArticleRepository

	produce events.Producer

	//v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	l      logger.Logger
}

func (a artService) GetPublishedById(ctx *gin.Context, id, uid int64) (domain.Article, error) {
	art, err := a.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			er := a.produce.ProduceReadEvent(ctx, events.ReadEvent{
				Uid: uid,
				Aid: id,
			})
			if er != nil {
				a.l.Error("发送读者阅读消息失败",
					logger.Int64("Uid", uid),
					logger.Int64("Aid", id),
					logger.Error(err))
			}
		}()
	}
	return art, err
}

func (a artService) GetById(ctx *gin.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a artService) List(ctx *gin.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a artService) Withdraw(ctx context.Context, art domain.Article) error {
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a artService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	// 制作库
	//id,err:=a.repo.Create(ctx,art)
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)

}

func (a artService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}

	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("保存到线上库失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("保存到线上库失败（重试彻底失败）",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
		//接入告警系统 可以手动处理
		//切换异步 写入本地文件
	}

	return id, err
}

func NewArticleService(repo article.ArticleRepository,
	l logger.Logger,
	produce events.Producer) ArticleService {
	return &artService{
		repo:    repo,
		l:       l,
		produce: produce,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository,
	l logger.Logger,
	produce events.Producer) ArticleService {
	return &artService{
		author:  author,
		reader:  reader,
		l:       l,
		produce: produce,
	}
}

func (a artService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a artService) update(ctx context.Context, art domain.Article) error {
	return a.repo.Update(ctx, art)
}
