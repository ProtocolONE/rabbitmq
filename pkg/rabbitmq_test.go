package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRabbitMq_ReConnect_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	broker, ok := b.(*Broker)
	assert.True(t, ok)

	tp1 := time.NewTimer(time.Second * 1)
	tp2 := time.NewTimer(time.Second * 2)
	tpEx := time.NewTimer(time.Second * 3)
	exit := make(chan bool, 1)

	broker.rabbitMQ.amqpUrl = ""

	notifyClose := make(chan *amqp.Error)
	broker.rabbitMQ.conn.NotifyClose(notifyClose)

	go func() {
		notifyClose <- &amqp.Error{Code: 111, Reason: "Test"}
		assert.False(t, broker.rabbitMQ.connected)
	}()

	for {
		select {
		case <-tp1.C:
			broker.rabbitMQ.amqpUrl = defaultAmqpUrl
		case <-tp2.C:
			assert.True(t, broker.rabbitMQ.connected)
		case <-tpEx.C:
			broker.rabbitMQ.close <- true
			exit <- true
		case <-exit:
			return
		}
	}
}

func TestRabbitMq_GetOpts_Default(t *testing.T) {
	br := &Broker{address: defaultAmqpUrl}
	rmq := br.newRabbitMq()

	defOpt := defaultQueueOpts
	res := rmq.getOpts(nil, defOpt)

	assert.Equal(t, defOpt, res)
}
