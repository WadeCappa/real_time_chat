package chatkafkamanager

// these tests assume that we have a kafka instance currently running

import (
	"log"
	"testing"

	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/constants"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/publisher"
)

func TestPublishWithoutError(t *testing.T) {
	err := publisher.PublishChatMessageToChannel([]string{"localhost:9092"}, 0, "test message", 0)
	if err != nil {
		t.Errorf("Failed to publish message %v", err)
	}
}

func TestPublishAndReadMessage(t *testing.T) {
	urls := []string{"localhost:9092"}
	done := make(chan bool)
	const channelId int64 = 12
	go func() {
		err := consumer.WatchChannel(urls, channelId, func(e events.Event, m consumer.Metadata) error {
			name, err := events.GetName(e)
			if err != nil || name == nil {
				t.Errorf("failed to get name %v", err)
			}
			if *name != constants.NEW_CHAT_MESSAGE_EVENT_NAME {
				t.Errorf("failed to get the correct event")
			} else {
				log.Printf("got event %v", e)
				done <- true
			}
			return nil
		})
		if err != nil {
			t.Errorf("Failed to consumer from stream %v", err)
		}
	}()
	err := publisher.PublishChatMessageToChannel(urls, 0, "test message", channelId)
	if err != nil {
		t.Errorf("Failed to publish message %v", err)
	}

	finished := <-done
	if !finished {
		t.Errorf("flag returned as false")
	}
}
