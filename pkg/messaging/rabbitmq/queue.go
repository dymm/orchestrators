package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

//Queue is a RabbitMQ queue
type Queue struct {
	name               string
	replyTo            string
	prefetchCount      int
	conn               *amqp.Connection
	connChannel        *amqp.Channel
	incommingFromQueue <-chan amqp.Delivery
	chanReturn         <-chan amqp.Return
	requestChan        chan messaging.WorkItem
}

//Close the connexion
func (i *Queue) Close() error {
	if i.connChannel != nil {
		i.connChannel.Close()
	}
	if i.requestChan != nil {
		close(i.requestChan)
		i.requestChan = nil
	}
	return nil
}

//Send a message to the 'reply to' queue
func (i Queue) Send(destination string, message messaging.WorkItem) error {
	serializedDictionnary, err := json.Marshal(message.GetValues())
	if err != nil {
		return errors.Wrapf(err, "Unable to marshal data '%v'", message.GetValues())
	}
	fmt.Printf("Sent to queue '%s' from queue '%s' '%s'\n", destination, i.name, string(serializedDictionnary))
	return i.connChannel.Publish(
		"",          // use the default exchange
		destination, // routing key, e.g. our queue name
		true,        // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(serializedDictionnary),
		})
}

//Receive a message from the queue
func (i Queue) Receive() (messaging.WorkItem, error) {
	var err error
	var values map[string]string
	select {
	case amqpDelivery, more := <-i.incommingFromQueue:
		if more == false {
			err = errors.New("Closed")
		} else {
			fmt.Println("Received on queue ", i.name, string(amqpDelivery.Body))
			err = json.Unmarshal(amqpDelivery.Body, &values)
		}
		if err != nil {
			amqpDelivery.Reject(false)
		}
	case chReturn, _ := <-i.chanReturn:
		if chReturn.ReplyText == "NO_ROUTE" {
			err = errors.Errorf("No route to channel '%s' from channel '%s'", chReturn.RoutingKey, i.name)
		} else {
			err = errors.Errorf("Output channel '%s' error %s", i.name, chReturn.ReplyText)
		}
	}
	return messaging.NewWorkItem(values), err
}
