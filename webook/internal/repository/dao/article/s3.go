package article

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"test/webook/internal/domain"
)

var statusPrivate = domain.ArticleStatusPrivate.ToUint8()

type S3DAO struct {
	oss *s3.S3
	GORMArticleDAO
	bucket *string
}

func NewOssDAO(oss *s3.S3, db *gorm.DB) ArticleDAO {
	return &S3DAO{
		oss:    oss,
		bucket: ekit.ToPtr[string]("webook-52171314"),
		GORMArticleDAO: GORMArticleDAO{
			db: db,
		},
	}
}

func (o *S3DAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (o *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id = art.Id
	)
	err := o.db.Transaction(func(tx *gorm.DB) error {
		var err error
		//制作库
		txDao := NewGORMArticleDAO(tx)
		if id > 0 {
			err = txDao.UpdateById(ctx, art)
		} else {
			id, err = txDao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		now := time.Now().UnixMilli()
		art.Id = id
		art.Ctime = now
		art.Utime = now
		art.Content = ""
		//OnConflict 代表数据冲突了
		err = tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  art.Title,
				"status": art.Status,
				"utime":  art.Utime,
			}),
		}).Create(&art).Error
		// 会生成  INSERT XX ON DUPLICATE KEY UPDATE XXX
		return err
	})
	if err != nil {
		return 0, err
	}
	_, err = o.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      o.bucket,
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	if err != nil {
		return 0, err
	}
	return id, err
}
