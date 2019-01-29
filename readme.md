RabbitMQ publisher/subscriber implementation 
=============

## Installation

`go get github.com/ProtocolONE/rabbitmq`

## Usage

```go
package main

import (
	"fmt"
	"github.com/ProtocolONE/payone-repository/pkg/proto/billing"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"rabbitmq/pkg"
)

type st struct {
	key         string
	retryBroker *rabbitmq.Broker
}

func main() {
	s := &st{
		key:         "test3",
		retryBroker: rabbitmq.NewBroker("amqp://127.0.0.1:5672"),
	}
	s.retryBroker.Opts.QueueOpts.Args = amqp.Table{
		"x-dead-letter-exchange":    "test3",
		"x-message-ttl":             int32(1 * 1000),
		"x-dead-letter-routing-key": "*",
	}
	s.retryBroker.Opts.ExchangeOpts.Name = "test.timeout10"

	br := rabbitmq.NewBroker("amqp://127.0.0.1:5672")

	if err := br.RegisterSubscriber("test3", s.test); err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 10; i++ {
		var err error

		st := &billing.Name{
			En: fmt.Sprintf("test %d", i),
			Ru: fmt.Sprintf("тест %d", i),
		}

		err = br.Publish("test3", st)

		if err != nil {
			log.Printf("Message publishing failed. Message id: %d ===> %s\n", i, err.Error())
			continue
		}

		log.Printf("Message publish success: %s\n", st.En)
	}

	if err := br.Subscribe(); err != nil {
		log.Fatalln(err)
	}

	log.Println("test")
}

func (s *st) test(a *billing.Name, b amqp.Delivery) (err error) {
	rnd := rand.Intn(100)

	if rnd > 50 {
		_ = s.retryBroker.Publish(s.key, a)

		log.Printf("[x] Failed fn - test: %s\n", a.En)
		return
	}

	log.Println("Complete fn - test: " + a.En)
	return
}
```