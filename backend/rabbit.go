package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func consumeRabbitEvents(rabbit *amqp.Connection, killChannel chan bool) (<-chan amqp.Delivery, error) {
	readChannel, err := rabbit.Channel()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	queue, err := readChannel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		readChannel.Close()
		return nil, err
	}

	go func() {
		for {
			shouldKill := <-killChannel
			if shouldKill {
				readChannel.QueueDelete(queue.Name, false, false, false)
				readChannel.Close()
			}
		}
	}()

	if err := readChannel.QueueBind(queue.Name, "", "events", false, nil); err != nil {
		killChannel <- true
		return nil, err
	}

	msgs, err := readChannel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		killChannel <- true
		return nil, err
	}
	return msgs, nil
}
