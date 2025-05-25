package internal_endpoints

import (
	"context"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/channel-manager/ownership"
	"github.com/WadeCappa/real_time_chat/channel-manager/store"
	"github.com/jackc/pgx/v5"
)

func CanWatch(postgresUrl string, ctx context.Context, request *external_channel_manager.CanWatchRequest) (*external_channel_manager.CanWatchResponse, error) {
	dbErr, err := store.Call(postgresUrl, func(c *pgx.Conn) error {
		err := ownership.CheckChannelPublic(request.ChannelId, c)
		if err == nil {
			return nil
		}

		return ownership.CheckUserCanEditChannel(request.ChannelId, request.UserId, c)
	})

	if dbErr != nil {
		return nil, *dbErr
	}

	if err != nil {
		return nil, err
	}

	return &external_channel_manager.CanWatchResponse{}, nil
}
