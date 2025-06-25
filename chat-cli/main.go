package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/chat-db/chat_db"
	"github.com/WadeCappa/real_time_chat/chat-watcher/chat_watcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	NO_COMMAND         = "no-command"
	NO_TOKEN           = ""
	NO_ADDRESS         = ""
	DEFAULT_CHANNEL_ID = 0

	POST_COMMAND           = "post"
	WATCH_COMMAND          = "watch"
	GET_CHANNELS_COMMAND   = "get-channels"
	CREATE_CHANNEL_COMMAND = "create-channel"
)

var commands = map[string]func() error{
	POST_COMMAND:           post,
	WATCH_COMMAND:          watch,
	GET_CHANNELS_COMMAND:   getChannels,
	CREATE_CHANNEL_COMMAND: createChannel,
}

var (
	addr      = flag.String("addr", NO_ADDRESS, "the server to hit")
	cmd       = flag.String("cmd", NO_COMMAND, "choose one of the following; ")
	userToken = flag.String("token", NO_TOKEN, "the user token to perform an operation with")
	channelId = flag.Int64("channel", DEFAULT_CHANNEL_ID, "the channel id on which to operate")
)

func withConnection(consumer func(*grpc.ClientConn) error) error {
	creds := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Printf("looking at %s\n", conn.CanonicalTarget())

	return consumer(conn)
}

func post() error {
	return withConnection(func(cc *grpc.ClientConn) error {
		fmt.Print("Enter message: ")
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to get input: %v", err)
		}

		newMetadata := metadata.Pairs("Authorization", *userToken)
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := chat_db.NewChatdbClient(cc)
		response, err := c.PublishMessage(newContext, &chat_db.PublishMessageRequest{ChannelId: *channelId, Message: message})

		if err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}

		log.Printf("%v\n", response)

		return nil
	})
}

func getChannels() error {
	return withConnection(func(cc *grpc.ClientConn) error {
		newMetadata := metadata.Pairs("Authorization", *userToken)
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := external_channel_manager.NewExternalchannelmanagerClient(cc)
		response, err := c.GetChannels(newContext, &external_channel_manager.GetChannelsRequest{PrefixSearch: ""})

		if err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}

		for {
			e, err := response.Recv()
			if err != nil {
				return fmt.Errorf("failed to get next channel: %v", err)
			}

			log.Println(e)
		}
	})
}

func createChannel() error {
	return withConnection(func(cc *grpc.ClientConn) error {
		var name string
		fmt.Print("Enter channel name: ")
		fmt.Scanln(&name)

		newMetadata := metadata.Pairs("Authorization", *userToken)
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := external_channel_manager.NewExternalchannelmanagerClient(cc)
		response, err := c.CreateChannel(
			newContext,
			&external_channel_manager.CreateChannelRequest{
				Public: true,
				Name:   name})

		if err != nil {
			return fmt.Errorf("failed to create channel: %v", err)
		}

		log.Printf("created channel %d", response.ChannelId)
		return nil
	})
}

func watch() error {
	return withConnection(func(cc *grpc.ClientConn) error {
		newMetadata := metadata.Pairs("Authorization", *userToken)
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := chat_watcher.NewChatwatcherserverClient(cc)
		response, err := c.WatchChannel(newContext, &chat_watcher.WatchChannelRequest{ChannelId: *channelId})

		if err != nil {
			return fmt.Errorf("failed to watch channel: %v", err)
		}

		for {
			e, err := response.Recv()
			if err != nil {
				return fmt.Errorf("failed to get next event: %v", err)
			}

			log.Println(e)
		}
	})
}

func main() {
	flag.Parse()

	command := commands[*cmd]
	if command == nil {
		log.Fatalf("Invalid command %s", *cmd)
	}

	err := command()
	if err != nil {
		log.Fatal(err)
	}
}
