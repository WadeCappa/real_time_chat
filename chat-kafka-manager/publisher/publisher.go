package publisher

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/constants"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
)

func createMessage(e events.Event) ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	name, err := events.GetName(e)
	if err != nil || name == nil {
		return nil, err
	}

	eventAndHeader := events.EventAndHeader{
		Header: *name,
		Data:   data,
	}
	return json.Marshal(eventAndHeader)
}

func getPublisher(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// visible for testing
func PublishEvent(consumer func([]byte) (int64, error), event events.Event) (int64, error) {
	data, err := createMessage(event)
	if err != nil {
		return 0, err
	}

	return consumer(data)
}

func writeToTopic(brokersUrl []string, userId int64, message string, channelId int64, topic string) (int64, error) {
	newMessage := events.NewChatMessageEvent{
		Content:   message,
		ChannelId: channelId,
		UserId:    userId,
	}

	pub, err := getPublisher(brokersUrl)
	if err != nil {
		return 0, err
	}

	// kafka specific
	consumer := func(data []byte) (int64, error) {
		kafkaEvent := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(data),
		}

		partition, offset, err := pub.SendMessage(kafkaEvent)
		if err != nil {
			return 0, err
		}

		logMessage, err := json.Marshal(map[string]interface{}{
			"message":   "Published the following message",
			"data":      string(data),
			"partition": partition,
			"offset":    offset,
		})
		if err == nil {
			log.Print(string(logMessage))
		} else {
			log.Printf("Failed to log success! %v", err)
		}

		return offset, nil
	}

	return PublishEvent(consumer, newMessage)
}

func PublishChatMessageToChannel(brokersUrl []string, userId int64, message string, channelId int64) (int64, error) {
	topic := constants.GetChannelTopic(channelId)
	return writeToTopic(brokersUrl, userId, message, channelId, topic)
}
