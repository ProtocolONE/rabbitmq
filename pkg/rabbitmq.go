package rabbitmq

import (
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

const (
	OptDurable    = "durable"
	OptAutoDelete = "autoDelete"
	OptExclusive  = "exclusive"
	OptNoWait     = "noWait"
	OptAutoAck    = "autoAck"
	OptNoLocal    = "noLocal"
	OptMandatory  = "mandatory"
	OptImmediate  = "immediate"
	OptInternal   = "internal"

	HeaderXDeath      = "x-death"
	HeaderRoutingKeys = "routing-keys"

	errorNicConnection = "connection not open"
	errorNilChannel    = "channel not open"
)

var defaultExchangeOpts = Opts{
	OptDurable:    true,
	OptAutoDelete: false,
	OptInternal:   false,
	OptNoWait:     false,
}

var defaultQueueOpts = Opts{
	OptDurable:    true,
	OptAutoDelete: false,
	OptExclusive:  false,
	OptNoWait:     false,
}

var defaultConsumeOpts = Opts{
	OptAutoAck:   false,
	OptExclusive: false,
	OptNoLocal:   false,
	OptNoWait:    false,
}

var defaultPublishOpts = Opts{
	OptMandatory: false,
	OptImmediate: false,
}

type Opts map[string]bool

type rabbitMq struct {
	amqpUrl string

	conn    *amqp.Connection
	channel *amqp.Channel

	uuid string

	wg sync.WaitGroup
	sync.Mutex
	connected      bool
	close          chan bool
	waitConnection chan struct{}
}

func (b *Broker) newRabbitMq() (rmq *rabbitMq) {
	rmq = &rabbitMq{
		amqpUrl:        b.address,
		close:          make(chan bool),
		waitConnection: make(chan struct{}),
	}
	id, err := uuid.NewRandom()

	if err != nil {
		rmq.uuid = ""
	} else {
		rmq.uuid = id.String()
	}

	close(rmq.waitConnection)
	return
}

func (r *rabbitMq) DeclareExchange(name, kind string, Opts Opts, args amqp.Table) error {
	pOpts := r.getOpts(Opts, defaultExchangeOpts)

	return r.channel.ExchangeDeclare(
		name,
		kind,
		pOpts[OptDurable],
		pOpts[OptAutoDelete],
		pOpts[OptInternal],
		pOpts[OptNoWait],
		args,
	)
}

func (r *rabbitMq) DeclareQueue(name string, Opts Opts, args amqp.Table) error {
	pOpts := r.getOpts(Opts, defaultQueueOpts)

	_, err := r.channel.QueueDeclare(
		name,
		pOpts[OptDurable],
		pOpts[OptAutoDelete],
		pOpts[OptExclusive],
		pOpts[OptNoWait],
		args,
	)

	return err
}

func (r *rabbitMq) QueueBind(queue, key, exchange string, noWait bool, args amqp.Table) error {
	return r.channel.QueueBind(queue, key, exchange, noWait, args)
}

func (r *rabbitMq) Consume(queue string, Opts Opts, args amqp.Table) (<-chan amqp.Delivery, error) {
	pOpts := r.getOpts(Opts, defaultConsumeOpts)

	return r.channel.Consume(
		queue,
		r.uuid,
		pOpts[OptAutoAck],
		pOpts[OptExclusive],
		pOpts[OptNoLocal],
		pOpts[OptNoWait],
		args,
	)
}

func (r *rabbitMq) Publish(exchange, key string, Opts Opts, message amqp.Publishing) error {
	pOpts := r.getOpts(Opts, defaultPublishOpts)
	return r.channel.Publish(exchange, key, pOpts[OptMandatory], pOpts[OptImmediate], message)
}

func (r *rabbitMq) getOpts(rec, def Opts) Opts {
	if rec == nil {
		return def
	}

	proc := make(Opts)

	for k, v := range def {
		val, ok := rec[k]

		if !ok {
			val = v
		}

		proc[k] = val
	}

	return proc
}

func (r *rabbitMq) tryConnect() (err error) {
	if r.conn, err = amqp.Dial(r.amqpUrl); err != nil {
		return
	}

	if r.channel, err = r.conn.Channel(); err != nil {
		return
	}

	return
}

func (r *rabbitMq) connect() (err error) {
	if err = r.tryConnect(); err != nil {
		return
	}

	r.Lock()
	r.connected = true
	r.Unlock()

	go r.reconnect()
	return
}

func (r *rabbitMq) reconnect() {
	var connect bool

	for {
		if connect {
			if err := r.tryConnect(); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			r.Lock()
			r.connected = true
			r.Unlock()

			close(r.waitConnection)
		}

		connect = true
		notifyClose := make(chan *amqp.Error)
		r.conn.NotifyClose(notifyClose)

		select {
		case <-notifyClose:
			r.Lock()
			r.connected = false
			r.waitConnection = make(chan struct{})
			r.Unlock()
		case <-r.close:
			return
		}
	}
}
