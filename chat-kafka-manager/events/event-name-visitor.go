package events

import (
	"fmt"

	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/constants"
)

type eventNameVisitor struct {
	EventVisitor

	result string
}

func (v *eventNameVisitor) VisitNewChatMessageEvent(e NewChatMessageEvent) error {
	v.result = constants.NEW_CHAT_MESSAGE_EVENT_NAME
	return nil
}

func GetName(e Event) (*string, error) {
	var v eventNameVisitor
	err := e.Visit(&v)
	if err != nil {
		return nil, fmt.Errorf("failed to get event name: %v", err)
	}
	return &v.result, nil
}
