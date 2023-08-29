package repository

import (
	"context"
	"database/sql"

	"test/webook/internal/domain"
	"test/webook/internal/repository/cache"
	"test/webook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository struct {
	u     *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(u *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		u:     u,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.u.Insert(ctx, r.domainToEntity(u))
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	return r.u.Update(ctx, r.domainToEntity(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.u.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.u.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	ue, err := r.u.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	ud := r.entityToDomain(ue)
	err = r.cache.Set(ctx, ud)
	return ud, err
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		NickName: u.NickName,
		Birth:    u.Birth,
		Synopsis: u.Synopsis,
		Ctime:    u.Ctime,
		Utime:    u.Utime,
	}
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		NickName: u.NickName,
		Birth:    u.Birth,
		Synopsis: u.Synopsis,
		Ctime:    u.Ctime,
		Utime:    u.Utime,
	}
}
