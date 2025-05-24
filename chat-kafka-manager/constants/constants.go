package constants

import (
	"fmt"
)

const (
	CHANNEL_ID_PARTITION_SUFFIX = "channel-"
	NEW_CHAT_MESSAGE_EVENT_NAME = "new-chat-message"
	UNKNOWN_EVENT_NAME          = "unknown-event"
)

func GetChannelTopic(channelId int64) string {
	return fmt.Sprintf(CHANNEL_ID_PARTITION_SUFFIX+"%d", channelId)
}
