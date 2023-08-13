package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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
	}
}
