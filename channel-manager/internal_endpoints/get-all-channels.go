package internal_endpoints

import (
	"context"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/channel-manager/store"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

func GetAllChannels(postgresUrl string, request *external_channel_manager.GetAllChannelsRequest, server grpc.ServerStreamingServer[external_channel_manager.GetAllChannelsResponse]) error {
	dbErr, err := store.Call(postgresUrl, func(c *pgx.Conn) error {
		rows, err := c.Query(
			context.Background(),
			"select id, name from channels")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var channelId int64
			var channelName string

			err = rows.Scan(&channelId, &channelName)
			if err != nil {
				return err
			}

			server.Send(
				&external_channel_manager.GetAllChannelsResponse{ChannelId: channelId})
		}

		return rows.Err()
	})

	if dbErr != nil {
		return *dbErr
	}

	if err != nil {
		return err
	}

	return nil
}
