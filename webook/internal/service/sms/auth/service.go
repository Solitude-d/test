package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"test/webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key string
}

// Send biz代表线下申请的业务方的token
func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var tc Claims
	//只要解析成功 那么说明token就是我发出去的
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
