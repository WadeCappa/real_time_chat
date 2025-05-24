package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/WadeCappa/real_time_chat/chat-writer/chat_writer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	NO_COMMAND             = "no-command"
	DEFAULT_SERVER_ADDRESS = "localhost:50052"

	NO_TOKEN           = ""
	POST_COMMAND       = "post"
	DEFAULT_CHANNEL_ID = 211
)

var commands = map[string]func(*grpc.ClientConn){
	POST_COMMAND: post,
}

var (
	serverAddress = flag.String("addr", DEFAULT_SERVER_ADDRESS, "the address to connect to")
	cmd           = flag.String("cmd", NO_COMMAND, "choose one of the following; ")
	userToken     = flag.String("token", NO_TOKEN, "the user token to perform an operation with")
	channelId     = flag.Int64("channel", DEFAULT_CHANNEL_ID, "the channel id on which to operate")
)

func main() {
	flag.Parse()

	command := commands[*cmd]
	if command == nil {
		log.Fatalf("Invalid command %s", *cmd)
	}

	conn, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Printf("looking at %s\n", conn.CanonicalTarget())

	command(conn)
}

func post(conn *grpc.ClientConn) {
	var message string
	fmt.Print("Enter message: ")
	fmt.Scanln(&message)

	newMetadata := metadata.Pairs("Authorization", *userToken)
	newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

	c := chat_writer.NewChatwriterserverClient(conn)
	response, err := c.PublishMessage(newContext, &chat_writer.PublishMessageRequest{ChannelId: *channelId, Message: message})

	if err != nil {
		log.Fatalf("Failed to send message! %v\n", err)
	}

	log.Printf("%v\n", response)
}
