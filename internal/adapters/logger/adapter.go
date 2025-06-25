package logger

import (
	"context"
	client "github.com/webitel/webitel-go-kit/infra/logger_client"
	"github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
)

var cl client.Publisher = &Adapter{}

type Adapter struct {
	channel rabbitmq.Publisher
}

func (a *Adapter) Publish(ctx context.Context, exchange string, routingKey string, body []byte) error {
	return a.channel.Publish(context.Background(), exchange, routingKey, body, nil)
}
func New(pub rabbitmq.Publisher) (*Adapter, error) {
	return &Adapter{channel: pub}, nil
}
