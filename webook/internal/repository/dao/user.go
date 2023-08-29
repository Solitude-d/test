package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱或手机号码冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (d *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := d.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicate
		}
	}
	return nil
}

func (d *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (d *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (d *UserDao) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("id = ?", id).Find(&u).Error
	return u, err
}

func (d *UserDao) Update(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	err := d.db.WithContext(ctx).Model(&u).Updates(&u).Error
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string

	NickName string
	Birth    string
	Synopsis string

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
