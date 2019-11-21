package messaging

//WorkItem represents the data exchanged between the elements constituting the workflow
type WorkItem struct {
	data map[string]string // dictionary of JSON serialized data to work with
}

//NewWorkItem create a new message
func NewWorkItem(values map[string]string) WorkItem {
	return WorkItem{
		data: values,
	}
}

//GetValues returns the message data
func (message WorkItem) GetValues() map[string]string {
	return message.data
}
