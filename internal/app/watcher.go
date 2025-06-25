package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	cfg "github.com/webitel/cases/config"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/webitel-go-kit/infra/fts_client"
	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
	"github.com/webitel/webitel-go-kit/pkg/watcher"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, exchange string, routingKey string, body []byte, headers amqp091.Table) error
}

type TriggerObserver[T any, V any] struct {
	id         string
	amqpBroker Publisher
	config     *cfg.TriggerWatcherConfig
	logger     *slog.Logger
	converter  func(T) (V, error)
}

func NewTriggerObserver[T any, V any](amqpBroker Publisher, config *cfg.TriggerWatcherConfig, conv func(T) (V, error), log *slog.Logger) (*TriggerObserver[T, V], error) {
	amqpObserver := &TriggerObserver[T, V]{
		amqpBroker: amqpBroker,
		config:     config,
		id:         "Trigger Watcher",
		logger:     log,
		converter:  conv,
	}
	return amqpObserver, nil
}

func (cao *TriggerObserver[T, V]) GetId() string {
	return cao.id
}

func (cao *TriggerObserver[T, V]) Update(et watcher.EventType, args map[string]any) error {
	var domainId int64
	obj, ok := args["obj"].(T)
	if !ok {
		return fmt.Errorf("could not convert to %d", obj)
	}

	session, ok := args["session"].(auth.Auther)
	if ok {
		domainId = session.GetDomainId()
	} else if domainId, ok = args["domain_id"].(int64); !ok {
		return fmt.Errorf("could not found domain id")
	}

	message, err := cao.converter(obj)
	if err != nil {
		return err
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Determine routing key prefix based on type of obj
	var objStr string
	switch any(obj).(type) {
	case *cases.Case:
		objStr = model.ScopeCases
	case *cases.CaseLink:
		objStr = model.BrokerScopeCaseLinks
	case *cases.CaseComment:
		objStr = model.ScopeCaseComments
	default:
		return fmt.Errorf("unsupported object type %T", obj)
	}

	routingKey := cao.getRoutingKeyByEventType("cases", objStr, et, domainId)
	cao.logger.Debug(fmt.Sprintf("Trying to publish message to %s", routingKey))

	if objStr == model.ScopeCaseComments || objStr == model.BrokerScopeCaseLinks {
		routingKey = cao.getRoutingKeyByEventType("cases", "case", et, domainId)
	}

	return cao.amqpBroker.Publish(context.Background(), cao.config.ExchangeName, routingKey, data, nil)
}

func (cao *TriggerObserver[T, V]) getRoutingKeyByEventType(
	service string,
	object string,
	eventType watcher.EventType,
	domainId int64,
) string {
	return fmt.Sprintf(
		"%s.%s.%s.%d",
		service,
		object,
		strings.Replace(cao.config.TopicName, "*", string(eventType), 1),
		domainId,
	)
}

type LoggerObserver struct {
	id      string
	logger  *wlogger.ObjectedLogger
	timeout time.Duration
}

func NewLoggerObserver(logger *wlogger.Logger, objclass string, timeout time.Duration) (*LoggerObserver, error) {
	objectedLogger, err := logger.GetObjectedLogger(objclass)
	if err != nil {
		return nil, err
	}
	return &LoggerObserver{
		id:      fmt.Sprintf("%s logger", objclass),
		logger:  objectedLogger,
		timeout: timeout,
	}, nil
}

func (l *LoggerObserver) GetId() string {
	return l.id
}

func (l *LoggerObserver) Update(et watcher.EventType, args map[string]any) error {
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
	case watcher.EventTypeCreate:
		tp = wlogger.CreateAction
	case watcher.EventTypeDelete:
		tp = wlogger.DeleteAction
	case watcher.EventTypeUpdate:
		tp = wlogger.UpdateAction
	default:
		return watcher.ErrUnknownType
	}
	message, err := wlogger.NewMessage(auth.GetUserId(), auth.GetUserIp(), tp, strconv.FormatInt(id, 10), args["obj"])
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), l.timeout)
	defer cancelFunc()
	_, err = l.logger.SendContext(ctx, auth.GetDomainId(), message)
	return err
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

func (l *FullTextSearchObserver[T, V]) Update(et watcher.EventType, args map[string]any) error {
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

	case watcher.EventTypeCreate:
		err = l.client.Create(auth.GetDomainId(), l.objclass, id, neededType)
	case watcher.EventTypeDelete:
		err = l.client.Delete(auth.GetDomainId(), l.objclass, id)
	case watcher.EventTypeUpdate:
		err = l.client.Update(auth.GetDomainId(), l.objclass, id, neededType)
	default:
		return watcher.ErrUnknownType
	}
	if err != nil {
		return err
	}
	return nil
}

func (l *FullTextSearchObserver[T, V]) getRoutingKeyByEventType(eventType watcher.EventType) string {
	return fmt.Sprintf("%s_case_key", string(eventType))
}
