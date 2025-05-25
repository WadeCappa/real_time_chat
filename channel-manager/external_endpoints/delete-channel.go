package external_endpoints

import (
	"context"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/channel-manager/ownership"
	"github.com/WadeCappa/real_time_chat/channel-manager/store"
	"github.com/jackc/pgx/v5"
)

func DeleteChannel(postgresUrl string, userId int64, ctx context.Context, request *external_channel_manager.DeleteChannelRequest) (*external_channel_manager.DeleteChannelResponse, error) {
	dbErr, err := store.Call(postgresUrl, func(c *pgx.Conn) error {
		err := ownership.CheckUserCanEditChannel(request.ChannelId, userId, c)
		if err != nil {
			return err
		}

		_, err = c.Exec(
			context.Background(),
			"delete from channels where channelId = $1 cascade",
			request.ChannelId)

		return err
	})

	if dbErr != nil {
		return nil, *dbErr
	}

	if err != nil {
		return nil, err
	}

	return &external_channel_manager.DeleteChannelResponse{}, nil
}
