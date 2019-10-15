package messaging

//Queue which can send and receive message
type Queue interface {

	//GetName give the queue name
	GetName() string

	//Receive a message from the queue
	Receive() (WorkItem, error)

	//Send a message to the destination queue
	Send(destination string, message WorkItem) error
}
