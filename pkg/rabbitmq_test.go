package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRabbitMq_ReConnect_Error(t *testing.T) {
	b, _ := NewBroker(defaultAmqpUrl)

	tp1 := time.NewTimer(time.Second * 1)
	tp2 := time.NewTimer(time.Second * 2)
	tpEx := time.NewTimer(time.Second * 3)
	exit := make(chan bool, 1)

	b.rabbitMQ.amqpUrl = ""

	notifyClose := make(chan *amqp.Error)
	b.rabbitMQ.conn.NotifyClose(notifyClose)

	go func() {
		notifyClose <- &amqp.Error{Code: 111, Reason: "Test"}
		assert.False(t, b.rabbitMQ.connected)
	}()

	for {
		select {
		case <-tp1.C:
			b.rabbitMQ.amqpUrl = defaultAmqpUrl
		case <-tp2.C:
			assert.True(t, b.rabbitMQ.connected)
		case <-tpEx.C:
			b.rabbitMQ.close <- true
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
