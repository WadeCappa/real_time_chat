package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/WadeCappa/real_time_chat/chat-watcher/chat_watcher"
	"github.com/WadeCappa/real_time_chat/chat-writer/chat_writer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	NO_COMMAND               = "no-command"
	DEFAULT_WRITE_SERVER_URL = "localhost:50052"
	DEFAULT_WATCH_SERVER_URL = "localhost:50053"

	NO_TOKEN           = ""
	DEFAULT_CHANNEL_ID = 211

	POST_COMMAND  = "post"
	WATCH_COMMAND = "watch"
)

var commands = map[string]func() error{
	POST_COMMAND:  post,
	WATCH_COMMAND: watch,
}

var (
	writeServerAddress = flag.String("write-server-url", DEFAULT_WRITE_SERVER_URL, "the address for a write server to communicate with")
	watchServerAddress = flag.String("watch-server-url", DEFAULT_WATCH_SERVER_URL, "the address for a watch server to communicate with")

	cmd       = flag.String("cmd", NO_COMMAND, "choose one of the following; ")
	userToken = flag.String("token", NO_TOKEN, "the user token to perform an operation with")
	channelId = flag.Int64("channel", DEFAULT_CHANNEL_ID, "the channel id on which to operate")
)

func withConnection(addr string, consumer func(*grpc.ClientConn) error) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Printf("looking at %s\n", conn.CanonicalTarget())

	return consumer(conn)
}

func post() error {
	return withConnection(*writeServerAddress, func(cc *grpc.ClientConn) error {
		var message string
		fmt.Print("Enter message: ")
		fmt.Scanln(&message)

		newMetadata := metadata.Pairs("Authorization", *userToken)
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := chat_writer.NewChatwriterserverClient(cc)
		response, err := c.PublishMessage(newContext, &chat_writer.PublishMessageRequest{ChannelId: *channelId, Message: message})

		if err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}

		log.Printf("%v\n", response)

		return nil
	})
}

func watch() error {
	return withConnection(*watchServerAddress, func(cc *grpc.ClientConn) error {
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
