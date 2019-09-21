package messaging

//WorkItem represents the data exchanged between the elements constituting the workflow
type WorkItem struct {
	data interface{} //Data to work with
}

//NewWorkItem create a new message
func NewWorkItem(value interface{}) WorkItem {
	return WorkItem{
		data: value,
	}
}

//GetData returns the message data
func (message WorkItem) GetData() interface{} {
	return message.data
}
