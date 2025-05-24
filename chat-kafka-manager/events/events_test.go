package events

import (
	"testing"

	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/constants"
)

func TestEventNameVisitor(t *testing.T) {
	newChatEvent := NewChatMessageEvent{}
	name := GetName(newChatEvent)
	if name != constants.NEW_CHAT_MESSAGE_EVENT_NAME {
		t.Errorf("failed to get the correct name")
	}
}
