package rabbitmq

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"log"
	"reflect"
	"sync"
	"time"
)

const (
	defaultExchangeKind = "topic"
	defaultQueueBindKey = "*"
	defaultContentType  = "application/protobuf"
)

type handler struct {
	method reflect.Value
	reqEl  reflect.Type
}

type subscriber struct {
	topic    string
	handlers []*handler
	fn       func(msg amqp.Delivery)

	mtx    sync.Mutex
	mayRun bool
	rabbit *rabbitMq

	opts struct {
		*QueueOpts
		*ExchangeOpts
		*QueueBindOpts
		*ConsumeOpts
	}

	ext map[string]bool
}

func (b *Broker) initSubscriber(topic string, rmq *rabbitMq) (subs *subscriber) {
	subs = &subscriber{
		topic:    topic,
		rabbit:   rmq,
		handlers: []*handler{},
		ext:      make(map[string]bool),
	}

	subs.opts.ExchangeOpts = b.Opts.ExchangeOpts
	subs.opts.QueueOpts = b.Opts.QueueOpts
	subs.opts.QueueBindOpts = b.Opts.QueueBindOpts
	subs.opts.ConsumeOpts = b.Opts.ConsumeOpts

	if subs.opts.ExchangeOpts.Name == "" {
		subs.opts.ExchangeOpts.Name = topic
	}

	if subs.opts.QueueOpts.Name == "" {
		subs.opts.QueueOpts.Name = subs.opts.ExchangeOpts.Name + ".queue"
	}

	return
}

func (s *subscriber) Subscribe() (err error) {
	if s.rabbit.conn == nil {
		return errors.New(errorNicConnection)
	}

	if s.rabbit.channel == nil {
		return errors.New(errorNilChannel)
	}

	fn := func(msg amqp.Delivery) {
		if msg.ContentType != defaultContentType {
			if s.opts.ConsumeOpts.Opts[OptAutoAck] == false {
				_ = msg.Nack(false, false)
			}

			return
		}

		for _, h := range s.handlers {
			st := reflect.New(h.reqEl).Interface().(proto.Message)
			err = proto.Unmarshal(msg.Body, st)

			if err != nil {
				if s.opts.ConsumeOpts.Opts[OptAutoAck] == false {
					_ = msg.Nack(false, false)
				}
				log.Println("[*] Unknown message type, message skipped ")
				continue
			}

			returnValues := h.method.Call([]reflect.Value{reflect.ValueOf(st), reflect.ValueOf(msg)})

			if err := returnValues[0].Interface(); err != nil {
				if s.opts.ConsumeOpts.Opts[OptAutoAck] == false {
					_ = msg.Nack(false, false)
				}
			} else {
				if s.opts.ConsumeOpts.Opts[OptAutoAck] == false {
					_ = msg.Ack(false)
				}
			}
		}
	}

	s.fn = fn
	s.mayRun = true

	go s.resubscribe()

	return
}

func (s *subscriber) resubscribe() {
	minResubscribeDelay := 100 * time.Millisecond
	maxResubscribeDelay := 30 * time.Second
	expFactor := time.Duration(2)
	reSubscribeDelay := minResubscribeDelay

	for {
		s.mtx.Lock()
		mayRun := s.mayRun
		s.mtx.Unlock()

		if !mayRun {
			return
		}

		select {
		case <-s.rabbit.close:
			return
		case <-s.rabbit.waitConnection:
		}

		s.rabbit.Lock()
		if !s.rabbit.connected {
			s.rabbit.Unlock()
			continue
		}

		sub, err := s.consume()
		s.rabbit.Unlock()

		switch err {
		case nil:
			reSubscribeDelay = minResubscribeDelay
		default:
			if reSubscribeDelay > maxResubscribeDelay {
				reSubscribeDelay = maxResubscribeDelay
			}
			time.Sleep(reSubscribeDelay)
			reSubscribeDelay *= expFactor
			continue
		}

		for d := range sub {
			s.rabbit.wg.Add(1)

			go func(d amqp.Delivery) {
				s.fn(d)
				s.rabbit.wg.Done()
			}(d)
		}
	}
}

func (s *subscriber) consume() (dls <-chan amqp.Delivery, err error) {
	err = s.rabbit.DeclareExchange(
		s.opts.ExchangeOpts.Name,
		s.opts.ExchangeOpts.Kind,
		s.opts.ExchangeOpts.Opts,
		s.opts.ExchangeOpts.Args,
	)

	if err != nil {
		return
	}

	err = s.rabbit.DeclareQueue(s.opts.QueueOpts.Name, s.opts.QueueOpts.Opts, s.opts.QueueOpts.Args)

	if err != nil {
		return
	}

	err = s.rabbit.QueueBind(
		s.opts.QueueOpts.Name,
		s.opts.QueueBindOpts.Key,
		s.opts.ExchangeOpts.Name,
		s.opts.QueueBindOpts.NoWait,
		s.opts.QueueBindOpts.Args,
	)

	if err != nil {
		return
	}

	dls, err = s.rabbit.Consume(s.opts.QueueOpts.Name, s.opts.ConsumeOpts.Opts, s.opts.ConsumeOpts.Args)

	if err != nil {
		return
	}

	return
}
