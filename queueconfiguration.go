package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/messaging/rabbitmq"
)

func createMessageQueueOrDie(receiveQueueName string) messaging.Queue {
	config := getChannelConfigurationFromEnvOrDie(receiveQueueName)
	queue, err := rabbitmq.New(config)
	if err != nil {
		fmt.Println("Error while creating the queue", err)
		os.Exit(1)
	}
	return queue
}

//getChannelConfigurationFromEnv get the channel configuration
func getChannelConfigurationFromEnvOrDie(receiveQueueName string) rabbitmq.ChannelConfig {

	defaultTTL, _ := time.ParseDuration(os.Getenv("RABBITMQ_DEFAULT_TTL"))
	inputQueueConfig := rabbitmq.ChannelConfig{
		Address:    os.Getenv("RABBITMQ_ADDR"),
		Name:       receiveQueueName,
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

	return inputQueueConfig
}
