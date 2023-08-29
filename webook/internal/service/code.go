package service

import (
	"context"
	"fmt"
	"math/rand"

	"test/webook/internal/repository"
	"test/webook/internal/service/sms"
)

const codeTemplId = "xxxxxxxx"

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTemplId, []string{code}, phone)
	return err
}

func (svc *CodeService) generateCode() string {
	//0-999999 之间 不包含1000000
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
