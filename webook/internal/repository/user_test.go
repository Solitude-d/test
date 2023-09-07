package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"test/webook/internal/domain"
	"test/webook/internal/repository/cache"
	cachemocks "test/webook/internal/repository/cache/mocks"
	"test/webook/internal/repository/dao"
	daomocks "test/webook/internal/repository/dao/mocks"
)

func TestCacheUserRepository_FindByID(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		ctx context.Context
		Id  int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindByID(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							Valid:  true,
							String: "123@qq.com",
						},
						Password: "password",
						Phone: sql.NullString{
							Valid:  true,
							String: "15381526629",
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)

				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "password",
					Phone:    "15381526629",
					Ctime:    now.UnixMilli(),
					Utime:    now.UnixMilli(),
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			Id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "password",
				Phone:    "15381526629",
				Ctime:    now.UnixMilli(),
				Utime:    now.UnixMilli(),
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "password",
						Phone:    "15381526629",
						Ctime:    now.UnixMilli(),
						Utime:    now.UnixMilli(),
					}, nil)

				d := daomocks.NewMockUserDao(ctrl)
				return d, c
			},
			ctx: context.Background(),
			Id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "password",
				Phone:    "15381526629",
				Ctime:    now.UnixMilli(),
				Utime:    now.UnixMilli(),
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，查询也不成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindByID(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("没查到"))
				return d, c
			},
			ctx:      context.Background(),
			Id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("没查到"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindByID(tc.ctx, tc.Id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
