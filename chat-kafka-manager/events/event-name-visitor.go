package events

import (
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/constants"
)

type eventNameVisitor struct {
	result string
}

func (v *eventNameVisitor) visitNewChatMessageEvent(e NewChatMessageEvent) {
	v.result = constants.NEW_CHAT_MESSAGE_EVENT_NAME
}

func GetName(e Event) string {
	var v eventNameVisitor
	e.Visit(&v)
	return v.result
}
