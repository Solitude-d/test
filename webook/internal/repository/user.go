package repository

import (
	"context"
	"test/webook/internal/domain"
	"test/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	u *dao.UserDao
}

func NewUserRepository(u *dao.UserDao) *UserRepository {
	return &UserRepository{
		u: u,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.u.Insert(ctx, domain.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	return r.u.Update(ctx, domain.User{
		Id:       u.Id,
		NickName: u.NickName,
		Birth:    u.Birth,
		Synopsis: u.Synopsis,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.u.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.u.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		NickName: u.NickName,
		Birth:    u.Birth,
		Synopsis: u.Synopsis,
		Email:    u.Email,
	}, nil
}
