package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	AtKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")
	RtKey = []byte("miyn8y9abnd7q4zkq2m73yw8tu9j5ixz")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	return nil
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("session 已经无效了")
	default:
		return nil
	}
}

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	sess := strings.Split(tokenHeader, " ")
	if len(sess) != 2 {
		return ""
	}
	return sess[1]
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
