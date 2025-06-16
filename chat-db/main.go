package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"slices"
	"time"

	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/chat-db/chat_db"
	"github.com/WadeCappa/real_time_chat/chat-db/result"
	"github.com/WadeCappa/real_time_chat/chat-db/store"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/publisher"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

const (
	DEFAULT_KAFKA_HOSTNAME           = "localhost:9092"
	DEFAULT_AUTH_HOSTNAME            = "localhost:50051"
	DEFAULT_POSTGRES_URL             = "postgres://postgres:pass@localhost:5432/chat_db"
	DEFAULT_CHANNEL_MANAGER_HOSTNAME = "localhost:50055"
	DEFAULT_PORT                     = 50052
	DEFAULT_LOAD_BATCH_SIZE          = 3
)

var (
	authHostname            = flag.String("auth-hostname", DEFAULT_AUTH_HOSTNAME, "the hostname for the auth service")
	kafkaHostname           = flag.String("kafka-hostname", DEFAULT_KAFKA_HOSTNAME, "the hostname for kafka")
	postgresUrl             = flag.String("postgres-hostname", DEFAULT_POSTGRES_URL, "the hostname for postgres")
	channelManangerHostname = flag.String("channel-manager-hostname", DEFAULT_CHANNEL_MANAGER_HOSTNAME, "the hostname for the channel manager")
	port                    = flag.Int("port", DEFAULT_PORT, "port for this service")
)

type chatDbServer struct {
	chat_db.ChatdbServer
}

type message struct {
	channelId   int64
	userId      int64
	messageId   int64
	time_posted time.Time
	content     string
}

func (s *chatDbServer) PublishMessage(
	ctx context.Context,
	request *chat_db.PublishMessageRequest) (*chat_db.PublishMessageResponse, error) {
	userId, err := auth.AuthenticateUser(ctx, *authHostname)
	if err != nil {
		return nil, err
	}

	if userId == nil {
		log.Println("did not return a valid user id")
		return nil, fmt.Errorf("returned an invalid userid")
	}

	log.Printf("looking at userid of %d\n", *userId)

	offset, err := publisher.PublishChatMessageToChannel(
		[]string{*kafkaHostname},
		*userId,
		request.Message,
		request.ChannelId)
	if err != nil {
		return nil, fmt.Errorf("failed to write to kafka: %v", err)
	}

	res := store.Call(*postgresUrl, func(c *pgx.Conn) result.Result[any] {
		tag, err := c.Exec(context.Background(),
			"insert into messages (user_id, message_id, channel_id, time_posted, content) values ($1, $2, $3, $4, $5)",
			userId,
			offset,
			request.ChannelId,
			time.Now(),
			request.Message)
		if err != nil {
			return result.Failed[any](fmt.Errorf("failed to create new message: %v", err))
		}
		fmt.Printf("tag from new channel request, %s\n", tag)

		return result.Result[any]{Result: nil, Err: nil}
	})
	if res.Err != nil {
		return nil, fmt.Errorf("failed to store new chat message event even after writing to kafka. This means we have lost a chat message: %v", res.Err)
	}

	return &chat_db.PublishMessageResponse{}, nil
}

func (s *chatDbServer) ReadMostRecent(
	request *chat_db.ReadMostRecentRequest,
	server grpc.ServerStreamingServer[chat_db.ReadMostRecentResponse]) error {

	result := store.Call(*postgresUrl, func(c *pgx.Conn) result.Result[[]chat_db.ReadMostRecentResponse] {

		rows, err := c.Query(
			context.Background(),
			"select user_id, message_id, channel_id, time_posted, content from messages where channel_id = $1 order by message_id desc limit $2",
			request.ChannelId,
			DEFAULT_LOAD_BATCH_SIZE)

		if err != nil {
			return result.Failed[[]chat_db.ReadMostRecentResponse](err)
		}
		defer rows.Close()

		messages := make([]chat_db.ReadMostRecentResponse, 0)
		for rows.Next() {
			var message message
			err := rows.Scan(
				&message.userId,
				&message.messageId,
				&message.channelId,
				&message.time_posted,
				&message.content)
			if err != nil {
				return result.Failed[[]chat_db.ReadMostRecentResponse](fmt.Errorf("failed to run scanner %v", err))
			}
			messages = append(messages, chat_db.ReadMostRecentResponse{
				Message:            message.content,
				UserId:             message.userId,
				ChannelId:          message.channelId,
				MessageId:          message.messageId,
				TimePostedUnixTime: message.time_posted.Unix(),
			})
		}
		if err := rows.Err(); err != nil {
			return result.Failed[[]chat_db.ReadMostRecentResponse](fmt.Errorf("some error happened while reading from db %v", err))
		}

		slices.Reverse(messages)
		return result.Success(messages)
	})

	if result.Err != nil {
		return fmt.Errorf("failed to load messages %v", result.Err)
	}

	for _, message := range *result.Result {
		err := server.Send(&message)
		if err != nil {
			return fmt.Errorf("failed to send message %v", err)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	chat_db.RegisterChatdbServer(server, &chatDbServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
