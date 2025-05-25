package external_endpoints

import (
	"context"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/channel-manager/store"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

func GetChannels(postgresUrl string, userId int64, request *external_channel_manager.GetChannelsRequest, server grpc.ServerStreamingServer[external_channel_manager.GetChannelsResponse]) error {
	dbErr, err := store.Call(postgresUrl, func(c *pgx.Conn) error {
		rows, err := c.Query(
			context.Background(),
			"select id, name from channels where public = True and name like '%' || $1 || '%'",
			request.PrefixSearch)
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
				&external_channel_manager.GetChannelsResponse{
					ChannelId:   channelId,
					ChannelName: channelName})
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
