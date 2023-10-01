package logger

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
)

type MiddleBuilder struct {
	allowReqBody *atomic.Bool
	allowResBody *atomic.Bool
	loggerFun    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddleBuilder {
	return &MiddleBuilder{
		loggerFun:    fn,
		allowReqBody: atomic.NewBool(false),
		allowResBody: atomic.NewBool(false),
	}
}

func (b *MiddleBuilder) AllowReqBody(ok bool) *MiddleBuilder {
	b.allowReqBody.Store(ok)
	return b
}

func (b *MiddleBuilder) AllowResBody(ok bool) *MiddleBuilder {
	b.allowResBody.Store(ok)
	return b
}

func (b *MiddleBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			//body, _ := io.ReadAll(ctx.Request.Body)
			body, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		if b.allowResBody.Load() {
			ctx.Writer = responseWrite{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			//al.Duration = time.Now().Sub(start)
			b.loggerFun(ctx, al)
		}()

		//执行业务逻辑
		ctx.Next()

		//b.loggerFun(ctx, al)
	}
}

type responseWrite struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWrite) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w responseWrite) Write(data []byte) (int, error) {
	if len(data) > 1024 {
		w.al.ResBody = string(data[:1024])
	} else {
		w.al.ResBody = string(data)
	}
	return w.ResponseWriter.Write(data)
}
func (w responseWrite) WriteString(data string) (int, error) {
	if len(data) > 1024 {
		w.al.ResBody = data[:1024]
	} else {
		w.al.ResBody = data
	}
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	Method   string
	Url      string
	ReqBody  string
	ResBody  string
	Duration string
	Status   int
}
