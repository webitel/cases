package app

import (
	context "context"
	broker "github.com/webitel/cases/rabbit"
)

type LoggerAdapter struct {
	publisher *broker.RabbitBroker
}

func (l *LoggerAdapter) Publish(ctx context.Context, routingKey string, body []byte, headers map[string]any) error {
	return l.publisher.Publish("logger", routingKey, body, headers)
}

func NewLoggerAdapter(rabbit *broker.RabbitBroker) *LoggerAdapter {
	return &LoggerAdapter{
		publisher: rabbit,
	}
}
