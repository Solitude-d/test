package saramax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"

	"test/webook/pkg/logger"
)

type BatchHandler[T any] struct {
	l             logger.Logger
	fn            func(msgs []*sarama.ConsumerMessage, ts []T) error
	batchDuration time.Duration
	batchSize     int
}

func NewBatchHandler[T any](l logger.Logger,
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{l: l, fn: fn, batchDuration: time.Second, batchSize: 10}
}

func (b BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgch := claim.Messages()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, b.batchSize)
		ts := make([]T, 0, b.batchSize)
		for i := 0; i < b.batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgch:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化消息失败",
						logger.Error(err),
						logger.Strings("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset))
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			b.l.Error("调用业务批量接口失败",
				logger.Error(err))
		}
		for _, msg := range msgs {
			//提交
			session.MarkMessage(msg, "")
		}
	}
}
