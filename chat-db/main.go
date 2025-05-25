package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/WadeCappa/real_time_chat/chat-db/chat_db"
	"github.com/WadeCappa/real_time_chat/chat-db/store"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/consumer"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/events"
	"github.com/gocql/gocql"
	"google.golang.org/grpc"
)

const (
	DEFAULT_KAFKA_HOSTNAME     = "localhost:9092"
	DEFAULT_CASSANDRA_HOSTNAME = "localhost:9042"
	DEFAULT_PORT               = 50054
	DEFAULT_LOAD_BATCH_SIZE    = 100
	TESTING_CHANNEL_ID         = 211
)

var (
	kafkaHostname     = flag.String("kafka-hostname", DEFAULT_KAFKA_HOSTNAME, "the hostname for kafka")
	cassandraHostname = flag.String("cassandra-hostname", DEFAULT_CASSANDRA_HOSTNAME, "the hostname for kafka")
	port              = flag.Int("port", DEFAULT_PORT, "port for this service")
)

type chatDbServer struct {
	chat_db.ChatdbServer
}

type message struct {
	channelId   int64
	userId      int64
	offset      int64
	time_posted time.Time
	content     string
}

type updateDataVisitor struct {
	events.EventVisitor

	metadata consumer.Metadata
}

func (v *updateDataVisitor) VisitNewChatMessageEvent(e events.NewChatMessageEvent) error {
	_, err := store.Call(*cassandraHostname, func(s *gocql.Session) (*bool, error) {
		err := s.Query("insert into posts_db.messages (userId, offset, channelId, time_posted, content) values (?, ?, ?, ?, ?)",
			e.UserId, v.metadata.Offset, e.ChannelId, v.metadata.TimePosted, e.Content).Exec()
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("failed to store new chat message event: %v", err)
	}
	return nil
}

func (s *chatDbServer) ReadMostRecent(request *chat_db.ReadMostRecentRequest, server grpc.ServerStreamingServer[chat_db.ReadMostRecentResponse]) error {

	_, err := store.Call(*cassandraHostname, func(s *gocql.Session) (*bool, error) {
		scanner := s.Query(`select userId, offset, channelId, time_posted, content from posts_db.messages where channelId = ? limit ?`,
			request.ChannelId, DEFAULT_LOAD_BATCH_SIZE).Iter().Scanner()

		for scanner.Next() {
			var message message
			err := scanner.Scan(&message)
			if err != nil {
				return nil, fmt.Errorf("failed to run scanner %v", err)
			}
			err = server.Send(&chat_db.ReadMostRecentResponse{
				Message:            message.content,
				UserId:             message.userId,
				ChannelId:          message.channelId,
				Offset:             message.offset,
				TimePostedUnixTime: message.time_posted.Unix(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to send message %v", err)
			}
		}
		// scanner.Err() closes the iterator, so scanner nor iter should be used afterwards.
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("some error happened while reading from db %v", err)
		}

		return nil, nil
	})

	return err
}

func listenAndWrite(channelId int64, kafkaUrl string, result chan error) {
	err := consumer.WatchChannel([]string{kafkaUrl}, channelId, func(e events.Event, m consumer.Metadata) error {
		v := updateDataVisitor{metadata: m}
		err := e.Visit(&v)
		if err != nil {
			return fmt.Errorf("failed to visit data event %v", err)
		}
		return nil
	})

	result <- err
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	res := make(chan error)

	go listenAndWrite(TESTING_CHANNEL_ID, *kafkaHostname, res)

	go func() {
		err := <-res
		if err != nil {
			log.Fatalf("failed to write to db: %v", err)
		}
	}()

	s := grpc.NewServer()
	chat_db.RegisterChatdbServer(s, &chatDbServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
