package localchannel

import (
	"github.com/dymm/gorchestrator/pkg/messaging"
)

//New return a local channel queue
func New() messaging.Queue {
	return Queue{
		internalChannel: make(chan messaging.WorkItem),
	}
}

//Queue which can send and receive message
type Queue struct {
	internalChannel chan messaging.WorkItem
}

//Receive a message from the queue
func (queue Queue) Receive() (messaging.WorkItem, error) {
	message := <-queue.internalChannel
	return message, nil
}

//Send a message into the queue
func (queue Queue) Send(message messaging.WorkItem) error {
	queue.internalChannel <- message
	return nil
}
