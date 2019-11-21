package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dymm/orchestrators/messageQ/pkg/messaging"
	"github.com/dymm/orchestrators/messageQ/pkg/messaging/rabbitmq"
)

//CreateMQMessageQueueOrDie get the communication channel
func CreateMQMessageQueueOrDie() messaging.Queue {

	defaultTTL, _ := time.ParseDuration(os.Getenv("RABBITMQ_DEFAULT_TTL"))
	inputQueueConfig := rabbitmq.ChannelConfig{
		Address:    os.Getenv("RABBITMQ_ADDR"),
		Name:       os.Getenv("RABBITMQ_OWN_NAME"),
		TTL:        defaultTTL,
		DeadLetter: os.Getenv("RABBITMQ_DEFAULT_DEAD_LETTER"),
	}

	if ttl, ok := os.LookupEnv("RABBITMQ_OWN_TTL"); ok {
		inputQueueConfig.TTL, _ = time.ParseDuration(ttl)
	}
	if val, ok := os.LookupEnv("RABBITMQ_QOS"); ok {
		if prefetchCount, errS := strconv.Atoi(val); errS != nil {
			inputQueueConfig.PrefetchCount = prefetchCount
		}
	}
	if val, ok := os.LookupEnv("RABBITMQ_OWN_DEAD_LETTER"); ok {
		inputQueueConfig.DeadLetter = val
	}

	queue, err := rabbitmq.New(inputQueueConfig)
	if err != nil {
		fmt.Println("Error while creating the queue.", err)
		os.Exit(1)
	}
	return queue
}
