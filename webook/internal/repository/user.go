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

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	Update(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CacheUserRepository struct {
	u     dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(u dao.UserDao, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		u:     u,
		cache: c,
	}
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.u.Insert(ctx, r.domainToEntity(u))
}

func (r *CacheUserRepository) Update(ctx context.Context, u domain.User) error {
	return r.u.Update(ctx, r.domainToEntity(u))
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.u.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.u.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
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

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
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

func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
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
