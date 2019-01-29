RabbitMQ publisher/subscriber implementation 
=============

## Installation

`go get github.com/ProtocolONE/rabbitmq`

## Usage

```go
package main

import (
	"fmt"
	"github.com/ProtocolONE/rabbitmq/internal/proto"
	"github.com/ProtocolONE/rabbitmq/pkg"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

func main() {
	br1, err := rabbitmq.NewBroker("amqp://127.0.0.1:5672")

	if err != nil {
		log.Fatalln(err)
	}

	topic := "test"

	br1.Opts.ExchangeOpts.Opts = rabbitmq.Opts{rabbitmq.OptAutoDelete: true}
	br1.Opts.QueueOpts.Opts = rabbitmq.Opts{rabbitmq.OptAutoDelete: true}

	br2, err := rabbitmq.NewBroker("amqp://127.0.0.1:5672")

	if err != nil {
		log.Fatalln(err)
	}

	br2.Opts.QueueOpts.Args = amqp.Table{
		"x-dead-letter-exchange":    topic,
		"x-message-ttl":             int32(1 * 1000),
		"x-dead-letter-routing-key": "*",
	}
	br2.Opts.ExchangeOpts.Name = "test.timeout10"

	fn := func(msg *test.One, d amqp.Delivery) (err error) {
		rnd := rand.Intn(100)
		rtr := int32(0)

		if v, ok := d.Headers["x-retry-count"]; ok {
			rtr = v.(int32)
		}

		if rnd > 50 {
			if rtr > 10 {
				log.Printf("[x] Max retries count is ended. Delete message : %s\n", msg.Value)
			} else {
				_ = br2.Publish(d.RoutingKey, msg, amqp.Table{"x-retry-count": rtr + 1})
			}

			log.Printf("[x] Retry (retry number %d) failed message: %s\n", rtr, msg.Value)
			return
		}

		log.Printf("Message successfully processed (retry number %d): %s", rtr, msg.Value)
		return
	}

	err = br1.RegisterSubscriber(topic, fn)

	for i := 0; i < 10; i++ {
		one := &test.One{Value: fmt.Sprintf("%s_%d", topic, i)}
		err = br1.Publish(topic, one, nil)

		if err == nil {
			log.Printf("Message published: %s\n", one.Value)
		}
	}

	err = br1.Subscribe(nil)

	if err != nil {
		log.Fatalln(err)
	}
}
```