package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"test/webook/internal/web/jwt"
	"test/webook/pkg/logger"
)

var L logger.Logger

// WrapReq 包裹所有请求 只要出问题就统一在这里写日志
func WrapReq[T any](fn func(ctx *gin.Context, req T) (Result, err error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			L.Error("处理业务逻辑出错",
				logger.Strings("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.Strings("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, Result{
			Data: res,
		})
	}
}

func WrapToken[C jwt.UserClaim](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, c)
		if err != nil {
			L.Error("处理业务逻辑出错",
				logger.Strings("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.Strings("route", ctx.FullPath()),
				logger.Error(err))
		}

		ctx.JSON(http.StatusOK, Result{
			Data: res,
		})
	}
}

func WrapBodyAndToken[Req any, C jwt.UserClaim](fn func(ctx *gin.Context, req Req, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			return
		}
		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, req, c)
		if err != nil {
			L.Error("处理业务逻辑出错",
				logger.Strings("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.Strings("route", ctx.FullPath()),
				logger.Error(err))
		}

		ctx.JSON(http.StatusOK, Result{
			Data: res,
		})
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
