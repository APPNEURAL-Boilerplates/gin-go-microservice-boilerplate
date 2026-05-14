package events

import (
	"context"
	"log/slog"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, payload any) error
}

type LogPublisher struct {
	logger *slog.Logger
}

func NewLogPublisher(logger *slog.Logger) *LogPublisher {
	return &LogPublisher{logger: logger}
}

func (p *LogPublisher) Publish(ctx context.Context, topic string, payload any) error {
	_ = ctx
	p.logger.Info("event published", "topic", topic, "payload", payload)
	return nil
}
