package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	ijwt "test/webook/internal/web/jwt"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}

}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path ...string) *LoginJWTMiddlewareBuilder {
	for _, v := range path {
		l.paths = append(l.paths, v)
	}
	return l
}

func (l *LoginJWTMiddlewareBuilder) Builder() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if path == c.Request.URL.Path {
				c.Next()
			}
		}
		//tokenHeader := c.GetHeader("Authorization")
		//if tokenHeader == "" {
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//sess := strings.Split(tokenHeader, " ")
		//if len(sess) != 2 {
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//tokenStr := sess[1]
		tokenStr := l.ExtractToken(c)
		claims := &ijwt.UserClaim{}
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

		if claims.UserAgent != c.Request.UserAgent() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(c, claims.Ssid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//使用长短token 不需要这个刷新机制
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString([]byte("xHd&^OrleeXM@Yq40gfww%8S%eND1*md"))
		//	if err != nil {
		//		println(err)
		//	}
		//	c.Header("x-jwt-token", tokenStr)
		//}
		c.Set("claim", claims)
	}
}
