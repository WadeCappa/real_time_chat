package syncer

import (
	"context"
	"fmt"

	"github.com/WadeCappa/real_time_chat/chat-db/result"
	"github.com/WadeCappa/real_time_chat/chat-db/store"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/jackc/pgx/v5"
)

type updateDataVisitor struct {
	events.EventVisitor

	metadata         consumer.Metadata
	postgresHostname string
}

func (v *updateDataVisitor) VisitNewChatMessageEvent(e events.NewChatMessageEvent) error {
	res := store.Call(v.postgresHostname, func(c *pgx.Conn) result.Result[any] {
		tag, err := c.Exec(context.Background(),
			"insert into messages (user_id, message_id, channel_id, time_posted, content) values ($1, $2, $3, $4, $5)",
			e.UserId,
			v.metadata.Offset,
			e.ChannelId,
			v.metadata.TimePosted,
			e.Content)
		if err != nil {
			return result.Failed[any](fmt.Errorf("failed to create new message: %v", err))
		}
		fmt.Printf("tag from new channel request, %s\n", tag)
		return result.Result[any]{Result: nil, Err: nil}
	})
	if res.Err != nil {
		return fmt.Errorf("failed to store new chat message event: %v", res.Err)
	}
	return nil
}
