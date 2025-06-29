package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/IBM/sarama"
	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/chat-db/chat_db"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/WadeCappa/real_time_chat/chat-watcher/chat_watcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	DEFAULT_PORT = 50053
)

var (
	port = flag.Int("port", DEFAULT_PORT, "port for this service")
)

func getChannelManagerHostname() string {
	return os.Getenv("CHANNEL_MANAGER_HOSTNAME")
}

func getChatDbHostname() string {
	return os.Getenv("CHAT_WRITER_HOSTNAME")
}

type chatWatcherServer struct {
	chat_watcher.ChatwatcherserverServer
}

type createChannelEventVisitor struct {
	events.EventVisitor

	e chat_watcher.ChannelEvent
}

func (v *createChannelEventVisitor) VisitNewChatMessageEvent(e events.NewChatMessageEvent) error {
	v.e = chat_watcher.ChannelEvent{
		EventUnion: &chat_watcher.ChannelEvent_NewMessage{
			NewMessage: &chat_watcher.NewMessageEvent{
				Conent:    e.Content,
				UserId:    e.UserId,
				ChannelId: e.ChannelId}}}
	return nil
}

func getRecentMessages(channelId int64, consumer func(*chat_db.ReadMostRecentResponse) error) (*int64, error) {
	conn, err := grpc.NewClient(getChatDbHostname(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, err
	}
	defer conn.Close()

	newContext := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
	c := chat_db.NewChatdbClient(conn)
	response, err := c.ReadMostRecent(newContext, &chat_db.ReadMostRecentRequest{ChannelId: channelId})
	if err != nil {
		return nil, fmt.Errorf("failed to get response from db service: %v", err)
	}

	var offset int64 = sarama.OffsetNewest

	for {
		e, err := response.Recv()
		if err == io.EOF {
			if offset != sarama.OffsetNewest {
				offset = offset + 1
			}
			return &offset, nil
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get next event: %v", err)
		}

		err = consumer(e)
		if err != nil {
			return nil, fmt.Errorf("failed to consume event: %v", err)
		}

		offset = e.MessageId
	}
}

func canUserWatchChannel(channelId, userId int64) error {
	conn, err := grpc.NewClient(getChannelManagerHostname(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	newContext := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	c := external_channel_manager.NewExternalchannelmanagerClient(conn)
	_, err = c.CanWatch(
		newContext,
		&external_channel_manager.CanWatchRequest{
			ChannelId: channelId,
			UserId:    userId})
	if err != nil {
		return fmt.Errorf("failed permission check: %v", err)
	}

	return nil
}

func (s *chatWatcherServer) WatchChannel(request *chat_watcher.WatchChannelRequest, server grpc.ServerStreamingServer[chat_watcher.WatchChannelResponse]) error {

	userId, err := auth.AuthenticateUser(server.Context())
	if err != nil {
		return fmt.Errorf("failed authenticaion: %v", err)
	}

	err = canUserWatchChannel(request.ChannelId, *userId)
	if err != nil {
		return fmt.Errorf("failed permission check: %v", err)
	}

	log.Printf("getting caught up %d\n", request.ChannelId)
	offset, err := getRecentMessages(request.ChannelId, func(rmrr *chat_db.ReadMostRecentResponse) error {
		e := chat_watcher.ChannelEvent{EventUnion: &chat_watcher.ChannelEvent_NewMessage{
			NewMessage: &chat_watcher.NewMessageEvent{
				Conent:    rmrr.Message,
				UserId:    rmrr.UserId,
				ChannelId: rmrr.ChannelId,
				MessageId: rmrr.MessageId}}}

		log.Println(&e)
		err = server.Send(&chat_watcher.WatchChannelResponse{Event: &e})
		if err != nil {
			return fmt.Errorf("failed to send event while getting caught up: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to get last offset: %v", err)
	}

	log.Printf("read up to offset %d\n", *offset)
	return consumer.WatchChannel(request.ChannelId, sarama.OffsetNewest, func(e events.Event, m consumer.Metadata) error {
		log.Println(e)
		v := createChannelEventVisitor{}
		v.e = chat_watcher.ChannelEvent{
			EventUnion:         &chat_watcher.ChannelEvent_UnknownEvent{UnknownEvent: &chat_watcher.UnknownEvent{Description: fmt.Sprintf("%v", e)}},
			TimePostedUnixTime: m.TimePosted.Unix(),
			Offest:             m.Offset,
		}
		err := e.Visit(&v)
		if err != nil {
			return fmt.Errorf("failed to visit event: %v", err)
		}
		v.e.Offest = m.Offset
		err = server.Send(&chat_watcher.WatchChannelResponse{Event: &v.e})
		if err != nil {
			return fmt.Errorf("failed to send event: %v", err)
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
