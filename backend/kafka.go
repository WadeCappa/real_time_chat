package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

var kafkaHostnames = []string{"kafka"}
var topic string = "basic_topic"

func StartPublisher() (chan []byte, error) {
	publisher, err := getPublisher(kafkaHostnames)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	eventChannel := make(chan []byte)

	go func() error {
		defer publisher.Close()
		data := <-eventChannel
		event := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(data),
		}
		partition, offset, err := publisher.SendMessage(event)
		if err != nil {
			return err
		}
		fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
		return nil
	}()

	return eventChannel, nil
}

func StartSubscriber(topic string) (<-chan []byte, error) {
	subscriber, err := getSubscriber(kafkaHostnames)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	readEvents := make(chan []byte)

	go func() {
		defer subscriber.Close()
		consumer, err := subscriber.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			panic(err)
		}
		defer consumer.Close()

		for message := range consumer.Messages() {
			readEvents <- message.Value
		}
	}()

	return readEvents, nil
}

func getPublisher(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getSubscriber(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// NewConsumer creates a new consumer using the given broker addresses and configuration
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
