package rabbit

const (
	ExchangeTypeTopic = "topic"
)

type ExchangeDeclareOptions struct {
	Args       map[string]any
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type QueueDeclareOptions struct {
	Args       map[string]any
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

type ExchangeDeclareOption func(*ExchangeDeclareOptions)

func ExchangeEnableDurable(o *ExchangeDeclareOptions) {
	o.Durable = true
}

func ExchangeEnableNoWait(o *ExchangeDeclareOptions) {
	o.NoWait = true
}

type QueueDeclareOption func(options *QueueDeclareOptions)
