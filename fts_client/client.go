package ftsclient

import (
	"github.com/gammazero/deque"
	"github.com/webitel/cases/rabbit"
	client "github.com/webitel/webitel-go-kit/fts_client"
)

const DefaultQueueSize = 500

var cl client.Publisher = &DefaultClient{}

type message struct {
	exchange string
	rk       string
	body     []byte
}

type DefaultClient struct {
	channel *rabbit.RabbitBroker
	queue   *deque.Deque[*message]
}

func (f *DefaultClient) Send(exchange string, rk string, body []byte) error {
	err := f.channel.Publish(exchange, rk, body, nil)
	if err != nil {
		// Add message to the queue
		f.queue.PushBack(&message{
			exchange: exchange,
			rk:       rk,
			body:     body,
		})
		return err
	}
	// Try to process the queue
	if f.queue.Len() > 0 {
		for el := f.queue.PopFront(); f.queue.Len() > 0; {
			err = f.channel.Publish(el.exchange, el.rk, el.body, nil)
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

func NewDefaultClient(rabbit *rabbit.RabbitBroker) (*client.Client, error) {
	q := &deque.Deque[*message]{}
	q.SetBaseCap(DefaultQueueSize)
	c := &DefaultClient{channel: rabbit, queue: q}
	return client.New(c), nil
}
