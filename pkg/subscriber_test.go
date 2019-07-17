package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscriber_Subscribe_Connection_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	topic := "test.subscriber"
	sub := b.initSubscriber(topic)

	err := sub.rabbit.conn.Close()
	assert.Nil(t, err)

	sub.rabbit.conn = nil
	err = sub.Subscribe()

	assert.NotNil(t, err)
	assert.Equal(t, errorNicConnection, err.Error())
}

func TestSubscriber_Subscribe_Channel_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)
	topic := "test.subscriber"
	sub := b.initSubscriber(topic)

	err := sub.rabbit.channel.Close()
	assert.Nil(t, err)

	sub.rabbit.channel = nil
	err = sub.Subscribe()

	assert.NotNil(t, err)
	assert.Equal(t, errorNilChannel, err.Error())
}

func TestSubscriber_DeclareExchange_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test"
	b.subscriber = b.initSubscriber(topic)

	_, err := b.subscriber.consume()

	assert.Nil(t, err)
	b.subscriber.opts.ExchangeOpts.Opts = Opts{OptAutoDelete: false}

	_, err = b.subscriber.consume()
	assert.NotNil(t, err)
}

func TestSubscriber_DeclareQueue_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test"
	b.subscriber = b.initSubscriber(topic)

	_, err := b.subscriber.consume()

	assert.Nil(t, err)
	b.subscriber.opts.QueueOpts.Opts = Opts{OptAutoDelete: false}

	_, err = b.subscriber.consume()
	assert.NotNil(t, err)
}

func TestSubscriber_QueueBind_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test"
	b.subscriber = b.initSubscriber(topic)

	_, err := b.subscriber.consume()

	assert.Nil(t, err)

	assert.Nil(t, err)
	b.Opts.QueueBindOpts.Args = amqp.Table{"test": int(3)}

	_, err = b.subscriber.consume()
	assert.NotNil(t, err)
}

func TestSubscriber_Consume_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	b.Opts.ExchangeOpts.Opts = Opts{OptAutoDelete: true}
	b.Opts.QueueOpts.Opts = Opts{OptAutoDelete: true}

	topic := "test"
	b.subscriber = b.initSubscriber(topic)

	_, err := b.subscriber.consume()

	assert.Nil(t, err)

	assert.Nil(t, err)
	b.Opts.ConsumeOpts.Args = amqp.Table{"test": int(3)}

	_, err = b.subscriber.consume()
	assert.NotNil(t, err)
}