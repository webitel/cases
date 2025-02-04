package ftsclient

import (
	"github.com/gammazero/deque"
	amqp "github.com/rabbitmq/amqp091-go"
	client "github.com/webitel/webitel-fts/pkg/client"
)

const DefaultQueueSize = 500

var cl client.Publisher = &FtsClient{}

type message struct {
	exchange string
	rk       string
	body     []byte
}

type FtsClient struct {
	channel *amqp.Channel
	queue   *deque.Deque[*message]
}

func (f *FtsClient) Send(exchange string, rk string, body []byte) error {
	if f.channel.IsClosed() {
		// noop
		return amqp.ErrClosed
	}
	err := f.channel.Publish(exchange, rk, false, false, amqp.Publishing{
		Body:    body,
		Headers: amqp.Table{"content-type": "application/json"},
	})
	if err != nil {
		// Add message to the queue
		f.queue.PushBack(&message{
			exchange: exchange,
			rk:       rk,
			body:     body,
		})
		return err
	}
	// Try to clear the queue
	if f.queue.Len() > 0 {
		for el := f.queue.PopFront(); f.queue.Len() > 0; {
			err = f.channel.Publish(el.exchange, el.rk, false, false, amqp.Publishing{
				Body:    el.body,
				Headers: amqp.Table{"content-type": "application/json"},
			})
			if err != nil {
				// error occured whil clearing the queue
				// push get back the element to the front
				f.queue.PushFront(el)
				return err
			}
		}
	}

	return nil
}

func NewFtsClient(rabbit *amqp.Channel) (*client.Client, error) {
	q := &deque.Deque[*message]{}
	q.SetBaseCap(DefaultQueueSize)
	c := &FtsClient{channel: rabbit, queue: q}
	return client.New(c), nil
}
