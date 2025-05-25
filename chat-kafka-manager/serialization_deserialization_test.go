package chatkafkamanager

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/publisher"
)

func TestNewMessage(t *testing.T) {
	event := events.NewChatMessageEvent{Content: "test-message", UserId: 20, ChannelId: 987}
	publisher.PublishEvent(func(b []byte) error {
		recoveredEvent, err := consumer.GetEvent(b)
		if err != nil {
			t.Errorf("failed to get event")
		}
		oldName, err := events.GetName(event)
		if err != nil {
			t.Errorf("could not get name for original event %v", err)
		}
		newName, err := events.GetName(recoveredEvent)
		if err != nil {
			t.Errorf("could not get name for new event %v", err)
		}
		if *oldName != *newName {
			t.Errorf("did not get the same event type")
		}
		if event != recoveredEvent {
			t.Errorf("Some data was different")
		}
		return nil
	}, event)
}

func TestUnknownEvent(t *testing.T) {
	data, err := json.Marshal(map[string]interface{}{"unknown-key": "some-data"})
	if err != nil {
		t.Error("Failed to marshal json")
	}
	unknownEvent := events.EventAndHeader{Header: "unknown-event-type", Data: data}

	unknownData, err := json.Marshal(unknownEvent)
	if err != nil {
		t.Error("Failed to marshal unkonwn event")
	}

	_, expectedError := consumer.GetEvent(unknownData)
	if expectedError == nil {
		t.Errorf("event was successfully deserialized. This is unexpected since this event name should not exist")
	}

	if !strings.HasPrefix(expectedError.Error(), consumer.COULD_NOT_FIND_EVENT_NAME_ERROR) {
		t.Errorf("Got the wrong error %v\n", expectedError)
	}
}
