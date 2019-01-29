package rabbitmq

import (
	"github.com/ProtocolONE/rabbitmq/internal/proto"
	"github.com/gogo/protobuf/proto"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	defaultAmqpUrl = "amqp://127.0.0.1:5672"
)

type TestStruct struct{}

func TestNewBroker(t *testing.T) {
	b, err := NewBroker(defaultAmqpUrl)

	assert.Nil(t, err)

	assert.Equal(t, defaultAmqpUrl, b.address, "they should be equal")
	assert.NotNil(t, b.rabbitMQ, "broker.rabbitMq shouldn't be nil")

	assert.Nil(t, b.subscriber, "broker.subscriber should be nil")
	assert.Nil(t, b.publisher, "broker.publisher should be nil")
	assert.NotNil(t, b.Opts, "broker.Opts shouldn't be nil")
}

func TestNewBroker_Fail(t *testing.T) {
	_, err := NewBroker("")
	assert.NotNil(t, err)
}

func TestBroker_RegisterSubscriber_Correct(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(msg *test.One, b amqp.Delivery) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)
	assert.Nil(t, err)
	assert.Len(t, b.subscriber.handlers, 1)
}

func TestBroker_RegisterSubscriber_HandlerNotFunc(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := "func"
	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "handler must have a function type", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerEmptyArgs(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func() error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "handler func must have two income argument", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerCountArgs(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(a string, b int, c bool) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "handler func must have two income argument", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerIncorrectSecondArg(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(msg proto.Message, b int) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "second argument of handler func must have a amqp.Delivery type", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerIncorrectFirstArg(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(msg string, b amqp.Delivery) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "first argument of handler func must be pointer to struct", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerIncorrectFirstArgType(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(msg *TestStruct, b amqp.Delivery) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "first argument of handler func must be instance of a proto.Message interface", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerIncorrectOut(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(msg *test.One, b amqp.Delivery) {}

	err := b.RegisterSubscriber("test", fn)

	assert.Error(t, err)
	assert.Equal(t, "handler func must have outcome argument", err.Error())
	assert.Len(t, b.subscriber.handlers, 0)
}

func TestBroker_RegisterSubscriber_HandlerDuplicate(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn := func(msg *test.One, b amqp.Delivery) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn)
	assert.Nil(t, err)

	err = b.RegisterSubscriber("test", fn)
	assert.Error(t, err)
	assert.Equal(t, "handler func already subscribed", err.Error())

	assert.Len(t, b.subscriber.handlers, 1)
}

func TestBroker_RegisterSubscriber_HandlerIncorrectFirstArgTypes(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn1 := func(msg *test.One, b amqp.Delivery) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn1)
	assert.Nil(t, err)

	fn2 := func(msg *test.Two, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fn2)
	assert.Error(t, err)
	assert.Equal(t, "first arguments for all handlers must have equal types", err.Error())

	assert.Len(t, b.subscriber.handlers, 1)
}

func TestBroker_RegisterSubscriber_HandlerMoreOneHandlersCorrect(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	fn1 := func(msg *test.One, b amqp.Delivery) error {
		return nil
	}

	err := b.RegisterSubscriber("test", fn1)
	assert.Nil(t, err)

	fn2 := func(msg *test.One, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fn2)
	assert.Nil(t, err)

	fn3 := func(msg *test.One, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fn3)
	assert.Nil(t, err)

	assert.Len(t, b.subscriber.handlers, 3)
}

func TestBroker_NewPublisher(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	topic := "test.publisher"
	pub := b.newPublisher(topic)

	assert.Equal(t, b.rabbitMQ, pub.rabbit)
	assert.Equal(t, b.Opts.PublishOpts, pub.opts.PublishOpts)
	assert.Equal(t, b.Opts.QueueBindOpts, pub.opts.QueueBindOpts)
	assert.Equal(t, b.Opts.ConsumeOpts, pub.opts.ConsumeOpts)

	assert.Equal(t, topic, pub.opts.ExchangeOpts.Name)
	assert.Equal(t, pub.opts.ExchangeOpts.Name+".queue", pub.opts.QueueOpts.Name)
}

func TestBroker_InitSubscriber(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	topic := "test.subscriber"
	sub := b.initSubscriber(topic)

	assert.Equal(t, b.rabbitMQ, sub.rabbit)
	assert.Equal(t, topic, sub.topic)
	assert.Len(t, sub.handlers, 0)
	assert.Len(t, sub.ext, 0)

	assert.Equal(t, b.Opts.QueueBindOpts, sub.opts.QueueBindOpts)
	assert.Equal(t, b.Opts.ConsumeOpts, sub.opts.ConsumeOpts)

	assert.Equal(t, topic, sub.opts.ExchangeOpts.Name)
	assert.Equal(t, sub.opts.ExchangeOpts.Name+".queue", sub.opts.QueueOpts.Name)
}

func TestBroker_Publish(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	assert.Nil(t, b.publisher)

	topic := "test.publisher"
	one := &test.One{Value: topic}

	err := b.Publish(topic, one, nil)

	assert.Nil(t, err)
	assert.NotNil(t, b.publisher)
}

func TestBroker_Publish_MarshalFail(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)
	topic := "test.publisher"

	err := b.Publish(topic, nil, nil)
	assert.NotNil(t, err)
	assert.Regexp(t, "Marshal", err.Error())
}

func TestBroker_Subscribe(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)
	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	assert.Nil(t, b.subscriber)

	topic := "test"
	fn := func(msg *test.One, b amqp.Delivery) error {
		assert.Equal(t, topic, msg.Value)
		return nil
	}

	err := b.RegisterSubscriber(topic, fn)
	assert.Nil(t, err)
	assert.Len(t, b.subscriber.handlers, 1)

	one := &test.One{Value: topic}
	err = b.Publish(topic, one, nil)
	assert.Nil(t, err)

	tp := time.NewTimer(time.Second * 3)
	done := make(chan bool, 1)
	exit := make(chan bool, 1)

	go func(done chan bool) {
		err = b.Subscribe(done)
	}(done)

	assert.Nil(t, err)

	select {
	case <-tp.C:
		done <- true
		exit <- true
	}
	<-exit
}
