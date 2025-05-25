package external_endpoints

import (
	"context"
	"fmt"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/channel-manager/ownership"
	"github.com/WadeCappa/real_time_chat/channel-manager/store"
	"github.com/jackc/pgx/v5"
)

type createChannelResult struct {
	err       error
	channelId int64
}

func CreateChannel(postgresUrl string, userId int64, ctx context.Context, request *external_channel_manager.CreateChannelRequest) (*external_channel_manager.CreateChannelResponse, error) {
	result, err := store.Call(postgresUrl, func(c *pgx.Conn) createChannelResult {
		var newChannelId int64
		err := c.QueryRow(context.Background(), "select nextval('channel_ids')").Scan(&newChannelId)
		if err != nil {
			return createChannelResult{err: fmt.Errorf("failed to get new channel id: %v", err)}
		}

		tag, err := c.Exec(context.Background(),
			"insert into channels (id, name, public) values ($1, $2, $3)",
			newChannelId,
			request.Name,
			request.Public)
		fmt.Printf("tag from new channel request, %s\n", tag)
		if err != nil {
			return createChannelResult{err: fmt.Errorf("failed to write channel: %v", err)}
		}

		ownership.AddUserToChannel(newChannelId, userId, c)
		if err != nil {
			return createChannelResult{err: fmt.Errorf("failed to add user to channel: %v", err)}
		}

		return createChannelResult{err: nil, channelId: newChannelId}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create new channel: %v", err)
	}

	if result.err != nil {
		return nil, fmt.Errorf("failed during db call to create channel: %v", err)
	}

	return &external_channel_manager.CreateChannelResponse{ChannelId: result.channelId}, nil
}
