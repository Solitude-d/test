package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path ...string) *LoginMiddlewareBuilder {
	for _, v := range path {
		l.paths = append(l.paths, v)
	}
	return l
}

func (l *LoginMiddlewareBuilder) Builder() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(c *gin.Context) {
		//if c.Request.URL.Path == "/users/login" ||
		//	c.Request.URL.Path == "/users/signup" {
		//	c.Next()
		//}
		for _, path := range l.paths {
			if path == c.Request.URL.Path {
				c.Next()
			}
		}
		sess := sessions.Default(c)
		userid := sess.Get("userId")
		if userid == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Next()
		}
		updateTime := sess.Get("update_time")
		sess.Set("userId", userid)
		sess.Options(sessions.Options{
			MaxAge: 60 * 30,
		})
		now := time.Now().UnixMilli()
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			c.Next()
		}
		updateVal, ok := updateTime.(int64)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			c.Next()
		}
		if now-updateVal > 60*1000 {
			sess.Set("update_time", now)
			sess.Save()
			c.Next()
		}
	}
}
