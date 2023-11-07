package article

import (
	"context"

	"test/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}
