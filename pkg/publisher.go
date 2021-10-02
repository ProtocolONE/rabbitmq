package rabbitmq

import "github.com/streadway/amqp"

type publisher struct {
	rabbit *rabbitMq

	opts struct {
		*QueueOpts
		*ExchangeOpts
		*QueueBindOpts
		*ConsumeOpts
		*PublishOpts
	}
	autoCreateQueue bool
}

func (b *Broker) newPublisher(topic string) (pub *publisher) {
	pub = &publisher{rabbit: b.rabbitMQ, autoCreateQueue: b.autoCreateQueue}

	pub.opts.PublishOpts = b.Opts.PublishOpts
	pub.opts.ExchangeOpts = b.Opts.ExchangeOpts

	if pub.opts.ExchangeOpts.Name == "" {
		pub.opts.ExchangeOpts.Name = topic
	}

	if pub.autoCreateQueue {
		pub.opts.QueueOpts = b.Opts.QueueOpts
		pub.opts.QueueBindOpts = b.Opts.QueueBindOpts
		pub.opts.ConsumeOpts = b.Opts.ConsumeOpts

		if pub.opts.QueueOpts.Name == "" {
			pub.opts.QueueOpts.Name = pub.opts.ExchangeOpts.Name + ".queue"
		}
	}

	return
}

func (p *publisher) publish(topic string, msg amqp.Publishing) (err error) {
	err = p.rabbit.DeclareExchange(
		p.opts.ExchangeOpts.Name,
		p.opts.ExchangeOpts.Kind,
		p.opts.ExchangeOpts.Opts,
		p.opts.ExchangeOpts.Args,
	)

	if err != nil {
		return
	}

	if p.autoCreateQueue {
		err = p.rabbit.DeclareQueue(p.opts.QueueOpts.Name, p.opts.QueueOpts.Opts, p.opts.QueueOpts.Args)

		if err != nil {
			return
		}

		err = p.rabbit.QueueBind(
			p.opts.QueueOpts.Name,
			p.opts.QueueBindOpts.Key,
			p.opts.ExchangeOpts.Name,
			p.opts.QueueBindOpts.NoWait,
			p.opts.QueueBindOpts.Args,
		)

		if err != nil {
			return
		}
	}

	return p.rabbit.Publish(p.opts.ExchangeOpts.Name, topic, p.opts.PublishOpts.Opts, msg)
}
