package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/IBM/sarama"
	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/chat-db/chat_db"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/WadeCappa/real_time_chat/chat-watcher/chat_watcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	DEFAULT_KAFKA_HOSTNAME   = "localhost:9092"
	DEFAULT_AUTH_HOSTNAME    = "localhost:50051"
	DEFAULT_CHAT_DB_HOSTNAME = "localhost:50054"
	DEFAULT_PORT             = 50053
)

var (
	authHostname   = flag.String("auth-hostname", DEFAULT_AUTH_HOSTNAME, "the hostname for the auth service")
	kafkaHostname  = flag.String("kafka-hostname", DEFAULT_KAFKA_HOSTNAME, "the hostname for kafka")
	chatDbHostname = flag.String("chat-db-hostname", DEFAULT_CHAT_DB_HOSTNAME, "the hostname for the chat db service")
	port           = flag.Int("port", DEFAULT_PORT, "port for this service")
)

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
	conn, err := grpc.NewClient(*chatDbHostname, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

		offset = e.Offset
	}
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

	log.Printf("getting caught up %d\n", request.ChannelId)
	offset, err := getRecentMessages(request.ChannelId, func(rmrr *chat_db.ReadMostRecentResponse) error {
		e := chat_watcher.ChannelEvent{EventUnion: &chat_watcher.ChannelEvent_NewMessage{
			NewMessage: &chat_watcher.NewMessageEvent{
				Conent:    rmrr.Message,
				UserId:    rmrr.UserId,
				ChannelId: rmrr.ChannelId,
				MessageId: rmrr.Offset}}}

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

	// Right now, this has a race condition where if a message is published in between someone starts listening,
	// loads the previous ealiest messages, then starts reading, we'll missed that message that was just posted.
	// We can fix this if we pass the message_id in a new event type (now we'll have pre-store and post-store types
	// which is just a good idea anyway for encapsulation purposes), so that the new control flow is start listening
	// to the newest message which will tell us its message_id, then we'll load everything (within a limit) before that
	// message from the db, then we'll resume listening from that first message. This way we don't miss any
	// information.
	log.Printf("read up to offset %d\n", *offset)
	return consumer.WatchChannel([]string{*kafkaHostname}, request.ChannelId, sarama.OffsetNewest, func(e events.Event, m consumer.Metadata) error {
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
