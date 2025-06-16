package syncer

import (
	"context"
	"fmt"
	"io"
	"iter"
	"log"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getChannels(channelManagerHostname string) (iter.Seq[int64], error) {
	conn, err := grpc.NewClient(channelManagerHostname, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	log.Println("looking at ", conn.CanonicalTarget())

	c := external_channel_manager.NewExternalchannelmanagerClient(conn)
	resp, err := c.GetAllChannels(
		context.Background(),
		&external_channel_manager.GetAllChannelsRequest{})
	if err != nil {
		return nil, fmt.Errorf("query failed! %v", err)
	}

	return func(yield func(int64) bool) {
		for {
			channel, err := resp.Recv()
			if err == io.EOF {
				return
			}

			if !yield(channel.ChannelId) {
				return
			}
		}
	}, nil
}
