package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	"test/webook/internal/repository"
	"test/webook/pkg/logger"
	"test/webook/pkg/saramax"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	l      logger.Logger
	repo   repository.InteractiveRepository
}

func NewInteractiveReadEventBatchConsumer(client sarama.Client,
	l logger.Logger,
	repo repository.InteractiveRepository) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (k *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewBatchHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventBatchConsumer) Consume(msgs []*sarama.ConsumerMessage, ts []ReadEvent) error {
	bizs := make([]string, 0, len(ts))
	ids := make([]int64, 0, len(ts))
	for _, event := range ts {
		bizs = append(bizs, "article")
		ids = append(ids, event.Aid)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := k.repo.IncrReadBatchCnt(ctx, ids, bizs)
	if err != nil {
		k.l.Error("批量增加阅读计数失败",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil
}
