package localchannel

import (
	"errors"
	"fmt"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//New return a local channel queue
func New(name string, avaliableQueues map[string]messaging.Queue) messaging.Queue {
	return Queue{
		name:             name,
		internalChannel:  make(chan messaging.WorkItem, 2),
		accessibleQueues: avaliableQueues,
	}
}

//Queue which can send and receive message
type Queue struct {
	name             string
	internalChannel  chan messaging.WorkItem
	accessibleQueues map[string]messaging.Queue
}

//Receive a message from the queue
func (queue Queue) Receive() (messaging.WorkItem, error) {
	var err error
	message, isOpened := <-queue.internalChannel
	if isOpened == false {
		err = errors.New("Message queue closed")
	}
	return message, err
}

//Send a message into the queue
func (queue Queue) Send(destination string, message messaging.WorkItem) error {
	if destination == queue.name {
		queue.internalChannel <- message
		return nil
	}

	destQueue, found := queue.accessibleQueues[destination]
	if found == false {
		return fmt.Errorf("No queue '%s' accessible from '%s'", destination, queue.name)
	}
	return destQueue.Send(destination, message)
}
