package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	"test/webook/internal/repository"
	"test/webook/pkg/logger"
	"test/webook/pkg/saramax"
)

type InteractiveReadEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	repo   repository.InteractiveRepository
}

func NewInteractiveReadEventConsumer(client sarama.Client,
	l logger.Logger,
	repo repository.InteractiveRepository) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (k *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "article", t.Aid)
}
