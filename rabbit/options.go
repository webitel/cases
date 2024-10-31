package rabbit

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

func ExchangeEnableAutoDelete(o *ExchangeDeclareOptions) {
	o.AutoDelete = true
}

func ExchangeEnableInternal(o *ExchangeDeclareOptions) {
	o.Internal = true
}

func ExchangeEnableNoWait(o *ExchangeDeclareOptions) {
	o.NoWait = true
}

type QueueDeclareOption func(options *QueueDeclareOptions)

func QueueEnableDurable(o *QueueDeclareOptions) {
	o.Durable = true
}

func QueueEnableAutoDelete(o *QueueDeclareOptions) {
	o.AutoDelete = true
}

func QueueEnableExclusive(o *QueueDeclareOptions) {
	o.Exclusive = true
}

func QueueEnableNoWait(o *QueueDeclareOptions) {
	o.NoWait = true
}
