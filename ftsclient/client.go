package ftsclient

import (
	amqp "github.com/rabbitmq/amqp091-go"
	client "github.com/webitel/webitel-fts/pkg/client"
)

var cl client.Publisher = &FtsClient{}

type FtsClient struct {
	channel *amqp.Channel
}

func (f *FtsClient) Send(exchange string, rk string, body []byte) error {
	return f.channel.Publish(exchange, rk, false, false, amqp.Publishing{
		Body:    body,
		Headers: amqp.Table{"content-type": "application/json"},
	})
}

func NewFtsClient(rabbit *amqp.Channel) (*client.Client, error) {
	c := &FtsClient{channel: rabbit}
	return client.New(c), nil
}
