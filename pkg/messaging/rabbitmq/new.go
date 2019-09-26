package rabbitmq

import (
	"errors"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/streadway/amqp"
)

//New create / get a queue from a queue name
func New(config ChannelConfig) (queue messaging.Queue, err error) {

	var connexion *amqp.Connection
	connexion, err = amqp.Dial(config.Address)

	var amqpChannel *amqp.Channel
	if err == nil {
		amqpChannel, err = connexion.Channel()
	}
	if err != nil {
		return queue, errors.New("Input channel error. " + err.Error())
	}

	args := amqp.Table{}
	if config.DeadLetter != "" {
		args["x-dead-letter-exchange"] = config.DeadLetter + ".exchange"
	}
	if config.TTL.Seconds() > 0 {
		args["x-message-ttl"] = int32(config.TTL.Seconds() * 1000)
	}

	var q amqp.Queue
	q, err = amqpChannel.QueueDeclare(
		config.Name,       // name
		false,             // durable
		config.Name == "", // delete when unused
		false,             // exclusive
		false,             // no-wait
		args,              // arguments
	)

	if err != nil {
		return queue, errors.New("Channel error. " + err.Error())
	}

	//These will block all comm on the connexion until fully read
	incommingFromQueue, err := amqpChannel.Consume(
		config.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return queue, err
	}

	chanReturn := amqpChannel.NotifyReturn(make(chan amqp.Return))

	queue = &Queue{
		name:               q.Name,
		replyTo:            config.ReplyTo,
		conn:               connexion,
		connChannel:        amqpChannel,
		incommingFromQueue: incommingFromQueue,
		chanReturn:         chanReturn,
		prefetchCount:      config.PrefetchCount,
	}
	return
}
