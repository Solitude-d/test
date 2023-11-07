package dao

import (
	"gorm.io/gorm"

	"test/webook/internal/repository/dao/article"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{})
}
