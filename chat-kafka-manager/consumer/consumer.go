package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/constants"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
)

const (
	COULD_NOT_FIND_EVENT_NAME_ERROR = "could not find name for event"
)

func getEventMapper(name string) (func([]byte) (events.Event, error), error) {
	switch name {
	case constants.NEW_CHAT_MESSAGE_EVENT_NAME:
		return func(data []byte) (events.Event, error) {
			var event events.NewChatMessageEvent
			err := json.Unmarshal(data, &event)
			return event, err
		}, nil
	default:
		return nil, fmt.Errorf(COULD_NOT_FIND_EVENT_NAME_ERROR+": %s", name)
	}
}

// visible for testing
func GetEvent(data []byte) (events.Event, error) {
	var eventAndHeader events.EventAndHeader
	err := json.Unmarshal(data, &eventAndHeader)
	if err != nil {
		return nil, err
	}

	eventMapper, err := getEventMapper(eventAndHeader.Header)
	if err != nil {
		return nil, err
	}

	event, err := eventMapper(eventAndHeader.Data)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func getSubscriber(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getMetadata(_ *sarama.ConsumerMessage) (Metadata, error) {
	return Metadata{}, nil
}

func WatchChannel(brokersUrl []string, channelId int64, eventConsumer func(events.Event, Metadata) error) error {
	subscriber, err := getSubscriber(brokersUrl)
	if err != nil {
		return err
	}
	defer subscriber.Close()

	// We assume only one partition per channel since we'll be DB constrained anyway.
	consumer, err := subscriber.ConsumePartition(constants.GetChannelTopic(channelId), 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer consumer.Close()

	for message := range consumer.Messages() {
		event, err := GetEvent(message.Value)
		if err != nil {
			return err
		}
		metadata, err := getMetadata(message)
		if err != nil {
			return err
		}
		err = eventConsumer(event, metadata)
		if err != nil {
			return err
		}
	}

	return nil
}
