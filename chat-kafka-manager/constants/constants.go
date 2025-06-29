package constants

import (
	"fmt"
	"os"
)

const (
	CHANNEL_ID_PARTITION_SUFFIX = "channel-"
	ALL_MESSAGES_TOPIC          = "all-messages"
	NEW_CHAT_MESSAGE_EVENT_NAME = "new-chat-message"
	UNKNOWN_EVENT_NAME          = "unknown-event"
)

func GetChannelTopic(channelId int64) string {
	return fmt.Sprintf(CHANNEL_ID_PARTITION_SUFFIX+"%d", channelId)
}

func GetAllMessagesTopic() string {
	return ALL_MESSAGES_TOPIC
}

func GetKafkaHostname() []string {
	return []string{os.Getenv("KAFKA_HOSTNAME")}
}
