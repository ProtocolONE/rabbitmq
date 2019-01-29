package rabbitmq

import (
	"github.com/gogo/protobuf/proto/proto3_proto"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultAmqpUrl = "amqp://127.0.0.1:5672"
)

type TestStruct struct{}

func TestNewBroker(t *testing.T) {
	b := NewBroker(defaultAmqpUrl)

	assert.Equal(t, defaultAmqpUrl, b.address, "they should be equal")
	assert.NotNil(t, b.rabbitMQ, "broker.rabbitMq shouldn't be nil")

	assert.Nil(t, b.subscriber, "broker.subscriber should be nil")
	assert.Nil(t, b.publisher, "broker.publisher should be nil")
	assert.NotNil(t, b.Opts, "broker.Opts shouldn't be nil")
}

func TestRegisterSubscriber(t *testing.T) {
	b := NewBroker(defaultAmqpUrl)

	sFn := "func"
	err := b.RegisterSubscriber("test", sFn)

	assert.Error(t, err)
	assert.Equal(t, "handler must have a function type", err.Error())

	fnEmptyArg := func() error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnEmptyArg)

	assert.Error(t, err)
	assert.Equal(t, "handler func must have two income argument", err.Error())

	fnTreeArg := func(a string, b int, c bool) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnTreeArg)

	assert.Error(t, err)
	assert.Equal(t, "handler func must have two income argument", err.Error())

	fnIncorrectSecondArg := func(msg proto.Message, b int) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnIncorrectSecondArg)

	assert.Error(t, err)
	assert.Equal(t, "second argument of handler func must have a amqp.Table type", err.Error())

	fnIncorrectFirstArg := func(msg string, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnIncorrectFirstArg)

	assert.Error(t, err)
	assert.Equal(t, "first argument of handler func must be pointer to struct", err.Error())

	fnIncorrectFirstArgType := func(msg *TestStruct, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnIncorrectFirstArgType)

	assert.Error(t, err)
	assert.Equal(t, "first argument of handler func must be instance of a proto.Message interface", err.Error())

	fnIncorrectOut := func(msg *proto3_proto.Message, b amqp.Delivery) {}

	err = b.RegisterSubscriber("test", fnIncorrectOut)

	assert.Error(t, err)
	assert.Equal(t, "handler func must have outcome argument", err.Error())

	fnCorrect := func(msg *proto3_proto.Message, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnCorrect)
	assert.Nil(t, err)

	err = b.RegisterSubscriber("test", fnCorrect)
	assert.Error(t, err)
	assert.Equal(t, "handler func already subscribed", err.Error())

	fnCorrect2 := func(msg *proto3_proto.Nested, b amqp.Delivery) error {
		return nil
	}

	err = b.RegisterSubscriber("test", fnCorrect2)
	assert.Error(t, err)
	assert.Equal(t, "first arguments for all handlers must have equal types", err.Error())
}
