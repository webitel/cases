package rabbit

import (
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	conf "github.com/webitel/cases/config"
	cerror "github.com/webitel/cases/internal/errors"
)

const MaxReconnectAttempts = 100

type RabbitBroker struct {
	config            *conf.RabbitConfig
	connection        *amqp.Connection
	channel           *amqp.Channel
	amqpCloseNotifier chan *amqp.Error
	consumers         map[string]*rabbitQueueConsumer
	emergencyStopper  chan<- cerror.AppError
	gracefulStopper   chan any
}

func BuildRabbit(config *conf.RabbitConfig, errChan chan<- cerror.AppError) (*RabbitBroker, cerror.AppError) {
	return &RabbitBroker{
		config:           config,
		emergencyStopper: errChan,
		consumers:        make(map[string]*rabbitQueueConsumer),
		gracefulStopper:  make(chan any),
	}, nil
}

// Start starts the channel between rabbitMQ server and this server
func (l *RabbitBroker) Start() cerror.AppError {
	var appErr cerror.AppError
	appErr = l.connect()
	if appErr != nil {
		return appErr
	}

	appErr = l.StartAllConsumers()
	if appErr != nil {
		return appErr
	}
	go stopHandler(l)
	return nil
}

// Stop stops all consumers and connections of rabbit
func (l *RabbitBroker) Stop() {
	// send to the gracefulStopper message that signalizes graceful stop (accepted in the stopHandler)
	l.gracefulStopper <- "graceful"
	l.StopAllConsumers()
	if l.channel != nil {
		l.channel.Close()
	}
	if l.connection != nil {
		l.connection.Close()
	}
	defer slog.Info(fmtBrokerLog("connection gracefully closed"))
}

var stopHandler = func(l *RabbitBroker) {
	for {
		select {
		case amqpErr := <-l.amqpCloseNotifier:
			slog.Error(fmtBrokerLog(fmt.Sprintf("connection lost %s", amqpErr.Reason)), slog.Int("code", amqpErr.Code))

			var (
				continueReconnection = true
				reconnectAttempts    int
			)

			for continueReconnection {
				if reconnectAttempts >= MaxReconnectAttempts { // if max reconnect attempts reached -- stop execution
					l.emergencyStopper <- cerror.NewInternalError("app.broker.stop_handler_routine.reconnect_attempts.reached_limit", "max reconnection attempts")
					return
				}
				reconnectErr := l.reconnect()
				if reconnectErr != nil {
					reconnectAttempts++
					slog.Error(fmtBrokerLog(reconnectErr.Error()), slog.Int("attempt", reconnectAttempts))
					time.Sleep(time.Second * 5)
				} else {
					continueReconnection = false
				}

			}
		case <-l.gracefulStopper:
			return
		}
	}
}

func (l *RabbitBroker) connect() cerror.AppError {
	conn, err := amqp.Dial(l.config.Url)
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.listen.server_connect.fail", fmtBrokerLog(err.Error()))
	}
	l.connection = conn
	channel, err := conn.Channel()
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.listen.channel_connect.fail", fmtBrokerLog(err.Error()))
	}
	l.channel = channel
	l.amqpCloseNotifier = l.channel.NotifyClose(make(chan *amqp.Error))

	err = channel.Qos(1, 0, false)
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.listen.qos.fail", fmtBrokerLog(err.Error()))
	}
	slog.Info(fmtBrokerLog("connection and amqp channel are opened"))
	return nil
}

func (l *RabbitBroker) reconnect() cerror.AppError {
	// try to create new connection channel
	slog.Debug(fmtBrokerLog("trying to reconnect"))
	err := l.connect()
	if err != nil {
		return err
	}
	for s, consumer := range l.consumers {
		// make a new delivery channel with new connection
		ch, err := l.Consume(s, consumer.name)
		if err != nil {
			return err
		}
		consumer.delivery = ch
		// start listen to the new delivery channel
		consumer.Start()
	}
	return nil
}

// StopAllConsumers stops all consumers if exist
func (l *RabbitBroker) StopAllConsumers() cerror.AppError {
	for _, consumer := range l.consumers {
		consumer.Stop()
	}
	return nil
}

// StartAllConsumers starts all consumers if exist
func (l *RabbitBroker) StartAllConsumers() cerror.AppError {
	for _, consumer := range l.consumers {
		err := consumer.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *RabbitBroker) ExchangeDeclare(exchangeName string, kind string, opts ...ExchangeDeclareOption) cerror.AppError {
	var decarationOptions ExchangeDeclareOptions
	for _, opt := range opts {
		opt(&decarationOptions)
	}

	err := l.channel.ExchangeDeclare(
		exchangeName,                 // name
		kind,                         // type
		decarationOptions.Durable,    // durable
		decarationOptions.AutoDelete, // auto-deleted
		decarationOptions.Internal,   // internal
		decarationOptions.NoWait,     // no-wait
		decarationOptions.Args,       // arguments
	)
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.exchange_declare.request.fail", err.Error())
	}
	slog.Info(fmtBrokerLog(fmt.Sprintf("[%s] exchange declared", exchangeName)), slog.String("name", exchangeName))
	return nil
}

func (l *RabbitBroker) QueueDeclare(queueName string, opts ...QueueDeclareOption) (string, cerror.AppError) {
	var declarationOptions QueueDeclareOptions
	for _, opt := range opts {
		opt(&declarationOptions)
	}

	_, err := l.channel.QueueDeclare(
		queueName,
		declarationOptions.Durable,
		declarationOptions.AutoDelete,
		declarationOptions.Exclusive,
		declarationOptions.NoWait,
		declarationOptions.Args,
	)
	if err != nil {
		return "", cerror.NewInternalError("rabbit.listener.queue_declare.request.fail", err.Error())
	}
	slog.Info(fmtBrokerLog(fmt.Sprintf("[%s] queue declared", queueName)), slog.String("name", queueName))
	return queueName, nil
}

func (l *RabbitBroker) QueueBind(exchangeName string, queueName string, routingKey string, noWait bool, args map[string]any) cerror.AppError {
	err := l.channel.QueueBind(queueName, routingKey, exchangeName, noWait, args)
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.queue_bind.request.fail", err.Error())
	}
	slog.Info(fmtBrokerLog(fmt.Sprintf("[%s]->(%s)->[%s] queue bind", exchangeName, routingKey, queueName)), slog.String("exchange", exchangeName), slog.String("routing", routingKey), slog.String("receiver", queueName))
	return nil
}

func (l *RabbitBroker) Consume(queueName string, consumerName string) (<-chan amqp.Delivery, cerror.AppError) {
	ch, err := l.channel.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, cerror.NewInternalError("rabbit.listener.consume.request.fail", err.Error())
	}
	slog.Info(fmtBrokerLog(fmt.Sprintf("[%s] queue started to consume", queueName)), slog.String("name", queueName))
	return ch, nil
}

func (l *RabbitBroker) ExchangeBind(destination string, key string, source string, noWait bool, args map[string]any) cerror.AppError {
	err := l.channel.ExchangeBind(destination, key, source, noWait, args)
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.exchange_bind.request.fail", err.Error())
	}
	slog.Info(fmtBrokerLog(fmt.Sprintf("[%s]->(%s)->[%s] exchange binded", source, key, destination)), slog.String("source", source), slog.String("routing", key), slog.String("destination", destination))
	return nil
}

func (l *RabbitBroker) QueueStartConsume(queueName string, consumerName string, handleFunc HandleFunc, handleTimeout time.Duration) cerror.AppError {
	// make a connection
	ch, err := l.Consume(queueName, consumerName)
	if err != nil {
		return err
	}
	// initialize handler
	queue, err := BuildRabbitQueueConsumer(ch, handleFunc, consumerName, handleTimeout)
	if err != nil {
		return err
	}
	// start new consumer
	queue.Start()

	// insert handler in the registry
	l.consumers[queueName] = queue

	return nil
}

func (l *RabbitBroker) QueueStopConsume(queueName string) cerror.AppError {
	if consumer, consumerFound := l.consumers[queueName]; consumerFound {
		consumer.Stop()
		delete(l.consumers, queueName)
	}
	return nil
}

func fmtBrokerLog(description string) string {
	return fmt.Sprintf("broker: %s", description)
}

func (l *RabbitBroker) Publish(
	exchange string, // Exchange name -- [Cases]
	routingKey string, // Routing key -- [api path]
	body []byte, // Message body
	userId string, // User ID
	t time.Time, // Message timestamp
) cerror.AppError {
	err := l.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			UserId:      userId,
			Timestamp:   t,
		},
	)
	if err != nil {
		return cerror.NewInternalError("rabbit.listener.publish.request.fail", fmtBrokerLog(err.Error()))
	}
	slog.Info(fmtBrokerLog("cases message published"), slog.String("exchange", exchange), slog.String("routingKey", routingKey))
	return nil
}

func (l *RabbitBroker) GetChannel() *amqp.Channel {
	return l.channel
}
