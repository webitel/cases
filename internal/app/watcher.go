package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/webitel/cases/auth"
	model "github.com/webitel/cases/config"
	cerr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/rabbit"
	wlogger "github.com/webitel/logger/pkg/client/v2"
	"github.com/webitel/webitel-go-kit/fts_client"
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
	GetArgs() map[string]any
}

type Observer interface {
	Update(EventType, []byte, map[string]any) error
	GetId() string
}

type Watcher interface {
	Attach(EventType, Observer)
	Detach(EventType, Observer)
	OnEvent(et EventType, entity WatchMarshaller) error
}

// WatcherManager manages a clustered storage of watchers
type WatcherManager interface {
	// AddWatcher adds watcher to cluster
	AddWatcher(clusterId string, watcher Watcher)

	// RemoveCluster removes full cluster
	RemoveCluster(clusterId string)

	GetCluster(clusterId string) []Watcher

	// Notify full cluster with event
	Notify(clusterId string, et EventType, data WatchMarshaller) error

	Enable()

	Disable()

	GetState() bool
}

type DefaultWatcherManager struct {
	clusters map[string][]Watcher
	state    bool
}

func NewDefaultWatcherManager(state bool) *DefaultWatcherManager {
	return &DefaultWatcherManager{clusters: make(map[string][]Watcher), state: state}
}

func (d *DefaultWatcherManager) AddWatcher(clusterId string, watcher Watcher) {
	d.clusters[clusterId] = append(d.clusters[clusterId], watcher)
}

func (d *DefaultWatcherManager) RemoveCluster(clusterId string) {
	delete(d.clusters, clusterId)
}

func (d *DefaultWatcherManager) GetCluster(clusterId string) []Watcher {
	return d.clusters[clusterId]
}

func (d *DefaultWatcherManager) Notify(clusterId string, et EventType, data WatchMarshaller) error {
	if !d.state {
		// noop
		return nil
	}
	cl := d.GetCluster(clusterId)
	if cl == nil {
		return nil
	}
	for _, watcher := range cl {
		err := watcher.OnEvent(et, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DefaultWatcherManager) GetState() bool {
	return d.state
}

func (d *DefaultWatcherManager) Enable() {
	d.state = true
}

func (d *DefaultWatcherManager) Disable() {
	d.state = false
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
		err = o.Update(et, data, entity.GetArgs())
		if err != nil {
			slog.Error(fmt.Sprintf("observer %s: %s", o.GetId(), err.Error()))
		}
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

type CaseAMQPObserver[T any, V any] struct {
	id         string
	amqpBroker AMQPBroker
	config     *model.WatcherConfig
	logger     *slog.Logger
	converter  func(T, map[string]any) (V, error)
}

func NewCaseAMQPObserver[T any, V any](amqpBroker AMQPBroker, config *model.WatcherConfig, conv func(T, map[string]any) (V, error), log *slog.Logger) (*CaseAMQPObserver[T, V], error) {
	// declare exchange
	opts := []rabbit.ExchangeDeclareOption{rabbit.ExchangeEnableDurable, rabbit.ExchangeEnableNoWait}

	if err := amqpBroker.ExchangeDeclare(config.ExchangeName, rabbit.ExchangeTypeTopic, opts...); err != nil {
		return nil, fmt.Errorf("could not create topic exchange %s: %w", config.ExchangeName, err)
	}

	amqpObserver := &CaseAMQPObserver[T, V]{
		amqpBroker: amqpBroker,
		config:     config,
		id:         "Case AMQP Watcher",
		logger:     log,
		converter:  conv,
	}
	return amqpObserver, nil
}

func (cao *CaseAMQPObserver[T, V]) GetId() string {
	return cao.id
}

func (cao *CaseAMQPObserver[T, V]) Update(et EventType, _ []byte, args map[string]any) error {
	obj, ok := args["obj"].(T)
	if !ok {
		return fmt.Errorf("could not convert to %d", obj)
	}

	message, err := cao.converter(obj, args)
	if err != nil {
		return err
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	routingKey := cao.getRoutingKeyByEventType(et)
	cao.logger.Debug(fmt.Sprintf("Trying to piublish message to %s", routingKey))
	return cao.amqpBroker.Publish(cao.config.ExchangeName, routingKey, data, "", time.Now())
}

func (cao *CaseAMQPObserver[T, V]) getRoutingKeyByEventType(eventType EventType) string {
	return strings.Replace(cao.config.TopicName, "*", string(eventType), 1)
}

type LoggerObserver struct {
	id      string
	logger  *wlogger.ObjectedLogger
	timeout time.Duration
}

func NewLoggerObserver(logger *wlogger.LoggerClient, objclass string, timeout time.Duration) (*LoggerObserver, error) {
	return &LoggerObserver{
		id:      fmt.Sprintf("%s logger", objclass),
		logger:  logger.GetObjectedLogger(objclass),
		timeout: timeout,
	}, nil
}

func (l *LoggerObserver) GetId() string {
	return l.id
}

func (l *LoggerObserver) Update(et EventType, data []byte, args map[string]any) error {
	auth, ok := args["session"].(auth.Auther)
	if !ok {
		return fmt.Errorf("could not get session auth")
	}
	id, ok := args["id"].(int64)
	if !ok {
		return fmt.Errorf("could not get id")
	}
	var tp wlogger.Action
	switch et {
	case EventTypeCreate:
		tp = wlogger.CreateAction
	case EventTypeDelete:
		tp = wlogger.DeleteAction
	case EventTypeUpdate:
		tp = wlogger.UpdateAction
	default:
		return ErrUnknownType
	}
	message, err := wlogger.NewMessage(auth.GetUserId(), "", tp, id, args["obj"])
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()
	return l.logger.SendContext(ctx, auth.GetDomainId(), message)
}

type FullTextSearchObserver[T any, V any] struct {
	id        string
	client    *fts_client.Client
	objclass  string
	converter func(T, map[string]any) (V, error)
}

func NewFullTextSearchObserver[T any, V any](client *fts_client.Client, objclass string, converter func(T, map[string]any) (V, error)) (*FullTextSearchObserver[T, V], error) {
	return &FullTextSearchObserver[T, V]{
		id:        fmt.Sprintf("%s fts", objclass),
		client:    client,
		objclass:  objclass,
		converter: converter,
	}, nil
}

func (l *FullTextSearchObserver[T, V]) GetId() string {
	return l.id
}

func (l *FullTextSearchObserver[T, V]) Update(et EventType, _ []byte, args map[string]any) error {
	auth, ok := args["session"].(auth.Auther)
	if !ok {
		return fmt.Errorf("could not get session auth")
	}
	id, ok := args["id"].(int64)
	if !ok {
		return fmt.Errorf("could not get id")
	}
	obj, ok := args["obj"].(T)
	if !ok {
		return fmt.Errorf("could not convert to %d", obj)
	}

	neededType, err := l.converter(obj, args)
	if err != nil {
		return err
	}
	switch et {

	case EventTypeCreate:
		err = l.client.Create(auth.GetDomainId(), l.objclass, id, neededType)
	case EventTypeDelete:
		err = l.client.Delete(auth.GetDomainId(), l.objclass, id)
	case EventTypeUpdate:
		err = l.client.Update(auth.GetDomainId(), l.objclass, id, neededType)
	default:
		return ErrUnknownType
	}
	if err != nil {
		return err
	}
	return nil
}

func (l *FullTextSearchObserver[T, V]) getRoutingKeyByEventType(eventType EventType) string {
	return fmt.Sprintf("%s_case_key", string(eventType))
}
