package rabbitmq

import (
	"github.com/ProtocolONE/rabbitmq/internal/proto"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublisher_DeclareExchange_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test.publisher"
	one := &test.One{Value: topic}

	err := b.Publish(topic, one, nil)

	assert.Nil(t, err)
	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: false}

	err = b.Publish(topic, one, nil)
	assert.NotNil(t, err)
}

func TestPublisher_DeclareQueue_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test.publisher"
	one := &test.One{Value: topic}

	err := b.Publish(topic, one, nil)

	assert.Nil(t, err)
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: false}

	err = b.Publish(topic, one, nil)
	assert.NotNil(t, err)
}

func TestPublisher_QueueBind_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test.publisher"
	one := &test.One{Value: topic}

	err := b.Publish(topic, one, nil)

	assert.Nil(t, err)
	b.Opts.QueueBindOpts.Args = amqp.Table{"test": int(3)}

	err = b.Publish(topic, one, nil)
	assert.NotNil(t, err)
}
