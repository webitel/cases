package fts

import (
	"context"
	"github.com/gammazero/deque"
	client "github.com/webitel/webitel-go-kit/infra/fts_client"
	"github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
)

const DefaultQueueSize = 500

var cl client.Publisher = &DefaultClient{}

type message struct {
	rk   string
	body []byte
}

type DefaultClient struct {
	channel rabbitmq.Publisher
	queue   *deque.Deque[*message]
}

func (f *DefaultClient) Send(exchange string, rk string, body []byte) error {
	err := f.channel.Publish(context.Background(), exchange, rk, body, nil)
	if err != nil {
		// Add message to the queue
		f.queue.PushBack(&message{
			rk:   rk,
			body: body,
		})
		return err
	}
	// Try to process the queue
	if f.queue.Len() > 0 {
		for el := f.queue.PopFront(); f.queue.Len() > 0; {
			err = f.channel.Publish(context.Background(), exchange, el.rk, el.body, nil)
			if err != nil {
				// error occurred while clearing the queue
				// push get back the element to the front
				f.queue.PushFront(el)
				return err
			}
		}
	}

	return nil
}

func NewDefaultClient(pub rabbitmq.Publisher) (*DefaultClient, error) {
	q := &deque.Deque[*message]{}
	q.SetBaseCap(DefaultQueueSize)
	return &DefaultClient{channel: pub, queue: q}, nil
}
