package rabbitmq

import (
	"time"
)

//ChannelConfig store the message queue the channel configuration.
//Name défini le nom de la file de message.
//TTL indique le temp de vie des messages arrivant dans la file de message.
//DeadLetter permet de definir la file de message où les messages depassant le TTL seront transférés.
//ReplyTo est le nom de la file de message utilisée lors d'un envoi de message
type ChannelConfig struct {
	Address       string
	Name          string
	TTL           time.Duration
	DeadLetter    string
	ReplyTo       string
	PrefetchCount int
}
