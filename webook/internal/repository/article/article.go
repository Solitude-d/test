package article

import (
	"context"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"

	"test/webook/internal/domain"
	"test/webook/internal/repository"
	"test/webook/internal/repository/cache"
	"test/webook/internal/repository/dao/article"
	"test/webook/pkg/logger"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error
	List(ctx *gin.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx *gin.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx *gin.Context, id int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	dao      article.ArticleDAO
	userRepo repository.UserRepository

	cache cache.ArticleCache

	l logger.Logger
}

func (c *CacheArticleRepository) GetPublishedById(ctx *gin.Context, id int64) (domain.Article, error) {
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	usr, err := c.userRepo.FindByID(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.NickName,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (c *CacheArticleRepository) GetById(ctx *gin.Context, id int64) (domain.Article, error) {
	res, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(res), nil
}

func (c *CacheArticleRepository) List(ctx *gin.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				c.preCache(ctx, data)
			}()
			return data[:limit], err
		}
	}
	res, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return c.toDomain(src)
	})
	go func() {
		err := c.cache.SetFirstPage(ctx, uid, data)
		c.l.Error("article 第一页缓存失败", logger.Error(err))
		c.preCache(ctx, data)
	}()

	return data, nil
}

func (c *CacheArticleRepository) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, id, uid, status.ToUint8())
}

func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	if err == nil {
		c.cache.DelFirstPage(ctx, art.Author.Id)
		c.cache.SetPub(ctx, art)
	}
	return id, err
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	go func() {
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.Insert(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	go func() {
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

func (c *CacheArticleRepository) preCache(ctx *gin.Context, data []domain.Article) {
	if len(data) > 0 && len(data[0].Content) < 1024*1024 {
		err := c.cache.Set(ctx, data[0])
		if err != nil {
			c.l.Error("缓存预加载失败", logger.Error(err))
		}
	}
}

func NewArticleRepository(d article.ArticleDAO, l logger.Logger) ArticleRepository {
	return &CacheArticleRepository{
		dao: d,
		l:   l,
	}
}
