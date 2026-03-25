package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"

	customreg "github.com/webitel/custom/registry"
	customstore "github.com/webitel/custom/store"
)

func (app *App) initCustom() error {
	// ensure connected !
	store := app.Store
	// _, err := store.Database()
	// if err != nil {
	// 	return err
	// }
	custom := store.Custom()
	// [storage] as a default custom types resolver !
	resolver := customreg.GlobalTypes
	resolver = resolver.WithResolver(
		customstore.CustomTypeResolver(custom),
	)
	customreg.GlobalTypes = resolver

	go subscribeCustomDatasetUpdates(app)

	return nil
}

func subscribeCustomDatasetUpdates(app *App) {
	log := slog.Default()

	for {
		err := runCustomDatasetSubscription(app)
		if err != nil {
			log.Error("[CUSTOM::EVENT] subscription failed, reconnecting...", "error", err)
		} else {
			log.Warn("[CUSTOM::EVENT] subscription disconnected, reconnecting...")
		}

		time.Sleep(5 * time.Second)
	}
}

func runCustomDatasetSubscription(app *App) error {
	config := app.config
	rabbit, err := app.rabbitConn.Channel(context.Background())
	if err != nil {
		return err
	}

	var (
		// err        error
		args       amqp091.Table
		autoAck    = true
		exchange   = "webitel"
		consumer   = "custom" // uuid.NewString()
		queueName  = fmt.Sprintf("%s-%s", consumer, config.Consul.Id)
		bindingKey = "custom.dataset.#"
		handler    = clusterCustomDatasetEventHandler
	)
	_, err = rabbit.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // autoDelete
		true,      // exclusive
		false,     // noWait
		args,      // args
	)
	if err != nil {
		return err
	}

	err = rabbit.QueueBind(
		queueName,  // name
		bindingKey, // key
		exchange,   // exchange
		false,      // noWait
		args,       // args
	)
	if err != nil {
		return err
	}

	deliveries, err := rabbit.Consume(
		queueName, // queue
		consumer,  // consumer
		autoAck,   // autoAck
		false,     // exclusive
		false,     // nolocal
		false,     // nowait
		nil,       // args
	)
	if err != nil {
		return err
	}

	for recv := range deliveries {
		// handle devilvery message
		handler(recv)
	}
	// disconnected !
	return nil
}

func clusterCustomDatasetEventHandler(recv amqp091.Delivery) (_ error) {
	// [layout]: "custom.dataset.{event}.{dc}.{name}"
	const (
		_ = iota // routeWordConstCustom = iota
		_        // routeWordConstDataset

		routeWordCustomEvent // [ create | update | delete ]
		routeWordDomainId    // int64
		routeWordDatasetName

		routeWordMax
	)

	log := slog.Default()
	topic := recv.RoutingKey
	route := strings.Split(topic, ".")

	if len(route) != routeWordMax {
		log.Warn("[CUSTOM::EVENT]", "error", "invalid routing key", "topic", topic)
		return
	}

	// extract: {event}
	event := strings.ToLower(route[routeWordCustomEvent])
	switch event {
	case "create", "update", "delete":
		// OK ; well-known ...
	default:
		// ERR ; invalid event type
	}

	// extract: {dc}
	dc, _ := strconv.ParseInt(route[routeWordDomainId], 10, 64)
	if dc < 1 {
		log.Warn("[CUSTOM::EVENT]",
			"topic", topic, "event", event,
			"error", "invalid domain reference",
		)
		return
	}

	// extract: {path}
	path, _ := recv.Headers["dataset.path"].(string)
	if path == "" {
		log.Warn("[CUSTOM::EVENT]",
			"topic", topic, "event", event, "dc", dc,
			"error", "missing dataset reference",
		)
		return
	}

	_ = customreg.Invalidate(dc, path)

	log.Info("[CUSTOM::EVENT]",
		"topic", topic, "event", event, "dc", dc, "path", path,
	)

	return
}
