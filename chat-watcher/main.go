package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	chat_watcher "github.com/WadeCappa/real_time_chat/chat-watcher/chat-watcher"
	"google.golang.org/grpc"
)

const (
	DEFAULT_KAFKA_HOSTNAME = "localhost:9092"
	DEFAULT_AUTH_HOSTNAME  = "localhost:50051"
	DEFAULT_PORT           = 50053
)

var (
	authHostname  = flag.String("auth-hostname", DEFAULT_AUTH_HOSTNAME, "the hostname for the auth service")
	kafkaHostname = flag.String("kafka-hostname", DEFAULT_KAFKA_HOSTNAME, "the hostname for kafka")
	port          = flag.Int("port", DEFAULT_PORT, "port for this service")
)

type chatWatcherServer struct {
	chat_watcher.ChatwatcherserverServer
}

func (s *chatWatcherServer) WatchChannel(request *chat_watcher.WatchChannelRequest, server grpc.ServerStreamingServer[chat_watcher.WatchChannelResponse]) error {

	userId, err := auth.AuthenticateUser(server.Context(), *authHostname)
	if err != nil {
		return err
	}

	if userId == nil {
		log.Println("did not return a valid user id")
		return fmt.Errorf("returned an invalid userid")
	}

	return consumer.WatchChannel([]string{*kafkaHostname}, request.ChannelId, func(e events.Event, m consumer.Metadata) error {
		err := server.Send(&chat_watcher.WatchChannelResponse{Event: &chat_watcher.ChannelEvent{EventUnion: &chat_watcher.ChannelEvent_UnknownEvent{UnknownEvent: &chat_watcher.UnknownEvent{Description: "some event"}}}})
		if err != nil {
			return err
		}

		return nil
	})
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chat_watcher.RegisterChatwatcherserverServer(s, &chatWatcherServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
