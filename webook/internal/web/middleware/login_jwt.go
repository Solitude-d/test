package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"test/webook/internal/web"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) Builder() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if path == c.Request.URL.Path {
				c.Next()
			}
		}
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sess := strings.Split(tokenHeader, " ")
		if len(sess) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := sess[1]
		claims := &web.UserClaim{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("xHd&^OrleeXM@Yq40gfww%8S%eND1*md"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("xHd&^OrleeXM@Yq40gfww%8S%eND1*md"))
			if err != nil {
				println(err)
			}
			c.Header("x-jwt-token", tokenStr)
		}
		c.Set("claim", claims)
	}
}
