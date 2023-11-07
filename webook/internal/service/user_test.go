package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"test/webook/internal/domain"
	"test/webook/internal/repository"
	repomocks "test/webook/internal/repository/mocks"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now().UnixMilli()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$.N8hKDglFrgvWnJLpdw.L.0/gIaAdT3MRY23kS3IDPsSvkIn7jXNu",
						Phone:    "15381526629",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello123",
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$.N8hKDglFrgvWnJLpdw.L.0/gIaAdT3MRY23kS3IDPsSvkIn7jXNu",
				Phone:    "15381526629",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "hello123",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock DB 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "hello123",
			wantUser: domain.User{},
			wantErr:  errors.New("mock DB 错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$.N8hKDglFrgvWnJLpdw.L.0/gIaAdT3MRY23kS3IDPsSvkIn7jXNu",
						Phone:    "15381526629",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello123111",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), nil)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}

}
