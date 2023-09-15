package retryable

import (
	"context"
	"errors"

	"test/webook/internal/service/sms"
)

type Service struct {
	svc      sms.Service
	retryMax int
}

// Send 重试的短信发送实现
func (s Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, biz, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, biz, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试失败！")
}
