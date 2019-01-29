package rabbitmq

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
)

type Broker struct {
	address  string
	rabbitMQ *rabbitMq

	subscriber *subscriber
	publisher  *publisher

	Opts *BrokerOpts
}

type BrokerOpts struct {
	*QueueOpts
	*ExchangeOpts
	*QueueBindOpts
	*ConsumeOpts
	*PublishOpts
}

type QueueOpts struct {
	Name string
	Opts Opts
	Args amqp.Table
}

type ExchangeOpts struct {
	Name string
	Kind string
	Opts Opts
	Args amqp.Table
}

type QueueBindOpts struct {
	Key    string
	NoWait bool
	Args   amqp.Table
}

type ConsumeOpts struct {
	Opts Opts
	Args amqp.Table
}

type PublishOpts struct {
	Opts Opts
}

func NewBroker(address string) (b *Broker) {
	b = &Broker{address: address}
	b.init()

	rmq := b.newRabbitMq()
	err := rmq.connect()

	if err != nil {
		log.Fatalf("[*] RabbitMq connection failed with error: %s", err)
	}

	b.rabbitMQ = rmq

	return
}

func (b *Broker) init() {
	b.Opts = &BrokerOpts{
		ExchangeOpts: &ExchangeOpts{
			Kind: defaultExchangeKind,
			Opts: defaultExchangeOpts,
			Args: nil,
		},
		QueueOpts: &QueueOpts{
			Opts: defaultQueueOpts,
			Args: nil,
		},
		QueueBindOpts: &QueueBindOpts{
			Key:    defaultQueueBindKey,
			NoWait: false,
			Args:   nil,
		},
		ConsumeOpts: &ConsumeOpts{
			Opts: defaultConsumeOpts,
			Args: nil,
		},
		PublishOpts: &PublishOpts{Opts: defaultPublishOpts},
	}
}

func (b *Broker) RegisterSubscriber(topic string, fn interface{}) error {
	if b.subscriber == nil {
		b.subscriber = b.initSubscriber(topic, b.rabbitMQ)
	}

	typ := reflect.TypeOf(fn)

	if typ.Kind() != reflect.Func {
		return errors.New("handler must have a function type")
	}

	refFn := reflect.ValueOf(fn)
	fnName := runtime.FuncForPC(refFn.Pointer()).Name()
	key := fnName + refFn.String()

	if _, ok := b.subscriber.ext[key]; ok {
		return errors.New("handler func already subscribed")
	}

	tNum := typ.NumIn()

	if tNum != 2 {
		return errors.New("handler func must have two income argument")
	}

	if typ.In(1).Kind() != reflect.Struct {
		return errors.New("second argument of handler func must have a amqp.Table type")
	}

	reqType := typ.In(0)

	if reqType.Kind() != reflect.Ptr {
		return errors.New("first argument of handler func must be pointer to struct")
	}

	_, ok := reflect.New(reqType.Elem()).Interface().(proto.Message)

	if !ok {
		return errors.New("first argument of handler func must be instance of a proto.Message interface")
	}

	tNum = typ.NumOut()

	if tNum != 1 {
		return errors.New("handler func must have outcome argument")
	}

	if len(b.subscriber.handlers) > 0 {
		if b.subscriber.handlers[0].reqEl != reqType.Elem() {
			return errors.New("first arguments for all handlers must have equal types")
		}
	}

	h := &handler{method: refFn, reqEl: reqType.Elem()}
	b.subscriber.handlers = append(b.subscriber.handlers, h)

	b.subscriber.ext[key] = true

	return nil
}

func (b *Broker) Subscribe() (err error) {
	err = b.subscriber.Subscribe()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	// wait on kill signal
	case <-ch:
	}

	return
}

func (b *Broker) Publish(topic string, msg proto.Message) (err error) {
	if b.publisher == nil {
		b.publisher = b.newPublisher(topic)
	}

	m := amqp.Publishing{
		ContentType: defaultContentType,
	}

	m.Body, err = proto.Marshal(msg)

	if err != nil {
		return
	}

	return b.publisher.publish(topic, m)
}
