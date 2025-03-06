package app

import (
	"errors"
	"fmt"
	model "github.com/webitel/cases/config"
	cerr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/rabbit"
	"log/slog"
	"strings"
	"time"
)

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeDelete EventType = "remove"
	EventTypeUpdate EventType = "update"
)

var ErrUnknownType = errors.New("unknown type")

type WatchMarshaller interface {
	Marshal() ([]byte, error)
}

type Observer interface {
	Update(EventType, []byte) error
	GetId() string
}

type Watcher interface {
	Attach(EventType, Observer)
	Detach(EventType, Observer)
	OnEvent(et EventType, entity WatchMarshaller) error
}

type DefaultWatcher struct {
	observers map[EventType][]Observer
}

func NewDefaultWatcher() *DefaultWatcher {
	return &DefaultWatcher{
		observers: make(map[EventType][]Observer),
	}
}

func (dw *DefaultWatcher) Attach(et EventType, o Observer) {
	dw.observers[et] = append(dw.observers[et], o)
}
func (dw *DefaultWatcher) Detach(et EventType, o Observer) {
	for i, v := range dw.observers[et] {
		if v.GetId() == o.GetId() {
			dw.observers[et] = append(dw.observers[et][:i], dw.observers[et][i+1:]...)
			break
		}
	}
}

func (dw *DefaultWatcher) Notify(et EventType, entity WatchMarshaller) error {
	data, err := entity.Marshal()
	if err != nil {
		return err
	}
	for _, o := range dw.observers[et] {
		_ = o.Update(et, data)
	}
	return nil
}

func (dw *DefaultWatcher) OnEvent(et EventType, entity WatchMarshaller) error {
	switch et {
	case EventTypeCreate:
		return dw.OnCreate(entity)
	case EventTypeDelete:
		return dw.OnDelete(entity)
	case EventTypeUpdate:
		return dw.OnUpdate(entity)
	default:
		return ErrUnknownType
	}
}

func (dw *DefaultWatcher) OnCreate(entity WatchMarshaller) error {
	return dw.Notify(EventTypeCreate, entity)
}
func (dw *DefaultWatcher) OnDelete(entity WatchMarshaller) error {
	return dw.Notify(EventTypeDelete, entity)
}
func (dw *DefaultWatcher) OnUpdate(entity WatchMarshaller) error {
	return dw.Notify(EventTypeUpdate, entity)
}

type AMQPBroker interface {
	QueueDeclare(queueName string, opts ...rabbit.QueueDeclareOption) (string, cerr.AppError)
	ExchangeDeclare(exchangeName string, kind string, opts ...rabbit.ExchangeDeclareOption) cerr.AppError
	QueueBind(exchangeName string, queueName string, routingKey string, noWait bool, args map[string]any) cerr.AppError
	Publish(exchange string, routingKey string, body []byte, userId string, t time.Time) cerr.AppError
}
type CaseAMQPObserver struct {
	id         string
	amqpBroker AMQPBroker
	config     *model.WatcherConfig
	logger     *slog.Logger
}

func NewCaseAMQPObserver(amqpBroker AMQPBroker, config *model.WatcherConfig, log *slog.Logger) (*CaseAMQPObserver, error) {

	// TODO :: refactor: use package constant
	queueMessagesTTL := func(o *rabbit.QueueDeclareOptions) {
		if o == nil {
			return
		}
		o.Args = map[string]any{
			"x-message-ttl": config.QueuesMessagesTTL,
		}
	}

	// declare queue
	if _, err := amqpBroker.QueueDeclare(config.QueueName, rabbit.QueueEnableDurable, queueMessagesTTL); err != nil {
		return nil, fmt.Errorf("could not create create queue %s: %w", config.QueueName, err)
	}

	// declare exchange
	if err := amqpBroker.ExchangeDeclare(config.ExchangeName, "topic", rabbit.ExchangeEnableDurable); err != nil {
		return nil, fmt.Errorf("could not create direct exchange %s: %w", config.ExchangeName, err)
	}

	// bind queue
	err := amqpBroker.QueueBind(config.ExchangeName, config.QueueName, config.TopicName, false, nil)
	if err != nil {
		return nil, fmt.Errorf("could not bind create create_queue: %w", err)
	}

	amqpObserver := &CaseAMQPObserver{
		amqpBroker: amqpBroker,
		config:     config,
		id:         "Case AMQP Watcher",
		logger:     log,
	}
	return amqpObserver, nil
}

func (cao *CaseAMQPObserver) GetId() string {
	return cao.id
}

func (cao *CaseAMQPObserver) Update(et EventType, data []byte) error {
	routingKey := cao.getRoutingKeyByEventType(et)
	cao.logger.Debug(fmt.Sprintf("Trying to piublish message %s to %s", string(data), routingKey))
	return cao.amqpBroker.Publish(cao.config.ExchangeName, routingKey, data, cao.config.AMQPUser, time.Now())
}

func (cao *CaseAMQPObserver) getRoutingKeyByEventType(eventType EventType) string {
	return strings.Replace(cao.config.TopicName, "*", string(eventType), 1)
}
