package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"test/webook/internal/domain"
	"test/webook/internal/service"
	ijwt "test/webook/internal/web/jwt"
	"test/webook/pkg/ginx/middlewares"
	"test/webook/pkg/logger"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc      service.ArticleService
	l        logger.Logger
	intersvc service.Interactive
	biz      string
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
		biz: "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("articles")
	g.POST("/edit", h.Edit)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/publish", h.Publish)

	g.POST("/list", middlewares.WrapBodyAndToken[ListReq, ijwt.UserClaim](h.list))
	g.GET("/detail/:id", middlewares.WrapToken[ijwt.UserClaim](h.Detail))

	pub := server.Group("/pub")
	pub.GET("/:id", h.pubDetail)
	pub.POST("/like", middlewares.WrapBodyAndToken[LikeReq, ijwt.UserClaim](h.Like))
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaim)
	if !ok {
		ctx.JSON(http.StatusOK,
			Result{
				Code: 5,
				Msg:  "系统错误",
			})
		h.l.Error("用户session信息不存在")
		return
	}
	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK,
			Result{
				Code: 5,
				Msg:  "系统错误",
			})
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: req.Id,
	})
}
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaim)
	if !ok {
		ctx.JSON(http.StatusOK,
			Result{
				Code: 5,
				Msg:  "系统错误",
			})
		h.l.Error("用户session信息不存在")
		return
	}
	//省略检测输入
	id, err := h.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK,
			Result{
				Code: 5,
				Msg:  "系统错误",
			})
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaim)
	if !ok {
		ctx.JSON(http.StatusOK,
			Result{
				Code: 5,
				Msg:  "系统错误",
			})
		h.l.Error("用户session信息不存在")
		return
	}
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK,
			Result{
				Code: 5,
				Msg:  "系统错误",
			})
		h.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) list(ctx *gin.Context, req ListReq, uc ijwt.UserClaim) (middlewares.Result, error) {
	res, err := h.svc.List(ctx, uc.Uid, req.Offset, req.Limit)
	if err != nil {
		return middlewares.Result{
			Code: 5,
			Msg:  " 系统错误",
		}, nil
	}
	return middlewares.Result{
		Data: slice.Map[domain.Article, ArticleVO](res, func(idx int, src domain.Article) ArticleVO {
			return ArticleVO{
				Title:    src.Title,
				Id:       src.Id,
				Abstract: src.Abstract(),
				Status:   src.Status.ToUint8(),
				Ctime:    src.Ctime.Format(time.DateTime),
				Utime:    src.Utime.Format(time.DateTime),
			}
		}),
	}, nil
}

func (h *ArticleHandler) Detail(ctx *gin.Context, uc ijwt.UserClaim) (middlewares.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return middlewares.Result{
			Code: 4,
			Msg:  "参数错误",
		}, nil
	}
	res, err := h.svc.GetById(ctx, id)
	if err != nil {
		return middlewares.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	if res.Id != uc.Uid {
		return middlewares.Result{
			Code: 4,
			Msg:  "输入错误",
		}, fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", res.Id)
	}
	return middlewares.Result{
		Data: ArticleVO{
			Id:      res.Id,
			Title:   res.Title,
			Status:  res.Status.ToUint8(),
			Content: res.Content,
			Ctime:   res.Ctime.Format(time.DateTime),
			Utime:   res.Utime.Format(time.DateTime),
		},
	}, nil
}

func (h *ArticleHandler) pubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
	}
	var art domain.Article
	var eg errgroup.Group
	uc := ctx.MustGet("users").(ijwt.UserClaim)
	eg.Go(func() error {
		art, err = h.svc.GetPublishedById(ctx, id, uc.Uid)
		return err
	})
	//art, err := h.svc.GetPublishedById(ctx, id)
	//if err != nil {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 5,
	//		Msg:  "系统错误",
	//	})
	//	h.l.Error("获取文章信息失败", logger.Error(err))
	//	return
	//}
	var inter domain.Interactive

	eg.Go(func() error {
		inter, err = h.intersvc.Get(ctx, biz, id, uc.Uid)
		if err != nil {
			h.l.Error("获取文章信息失败", logger.Error(err))
		}
		return nil
	})
	//inter, err := h.intersvc.Get(ctx, biz, id, uc.Uid)
	//if err != nil {
	//	//这里可以出现错误 保证用户能继续浏览主要内容 只记录日志
	//	//ctx.JSON(http.StatusOK, Result{
	//	//	Code: 5,
	//	//	Msg:  "系统错误",
	//	//})
	//	h.l.Error("获取文章信息失败", logger.Error(err))
	//}
	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	go func() {
		er := h.intersvc.IncrReadCnt(ctx, h.biz, art.Id)
		if er != nil {
			h.l.Error("增加阅读计数失败", logger.Error(er))
		}
	}()
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:         art.Id,
			Title:      art.Title,
			Status:     art.Status.ToUint8(),
			Content:    art.Content,
			Author:     art.Author.Name,
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			Liked:      inter.Liked,
			LikeCnt:    inter.LikeCnt,
			Collected:  inter.Collected,
			CollectCnt: inter.CollectCnt,
			ReadCnt:    inter.ReadCnt,
		},
	},
	)
}

func (h *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc ijwt.UserClaim) (middlewares.Result, error) {
	var err error
	if req.Like {
		err = h.intersvc.Like(ctx, h.biz, req.Id, uc.Uid)
	} else {
		err = h.intersvc.CancelLike(ctx, h.biz, req.Id, uc.Uid)
	}
	if err != nil {
		return middlewares.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return middlewares.Result{
		Msg: "OK",
	}, nil
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
