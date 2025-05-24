package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/publisher"
	"github.com/WadeCappa/real_time_chat/chat-writer/chat_writer"
	"google.golang.org/grpc"
)

const (
	DEFAULT_KAFKA_HOSTNAME = "localhost:9092"
	DEFAULT_AUTH_HOSTNAME  = "localhost:50051"
	DEFAULT_PORT           = 50052
)

var (
	authHostname  = flag.String("auth-hostname", DEFAULT_AUTH_HOSTNAME, "the hostname for the auth service")
	kafkaHostname = flag.String("kafka-hostname", DEFAULT_KAFKA_HOSTNAME, "the hostname for kafka")
	port          = flag.Int("port", DEFAULT_PORT, "port for this service")
)

type chatWriterServer struct {
	chat_writer.ChatwriterserverServer
}

func (s *chatWriterServer) PublishMessage(ctx context.Context, request *chat_writer.PublishMessageRequest) (*chat_writer.PublishMessageResponse, error) {

	userId, err := auth.AuthenticateUser(ctx, *authHostname)
	if err != nil {
		return nil, err
	}

	if userId == nil {
		log.Println("did not return a valid user id")
		return nil, fmt.Errorf("returned an invalid userid")
	}

	log.Printf("looking at userid of %d\n", *userId)

	if err := publisher.PublishChatMessageToChannel([]string{*kafkaHostname}, *userId, request.Message, request.ChannelId); err != nil {
		return nil, err
	}

	return &chat_writer.PublishMessageResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chat_writer.RegisterChatwriterserverServer(s, &chatWriterServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
