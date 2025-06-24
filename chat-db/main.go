package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"time"

	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/chat-db/chat_db"
	"github.com/WadeCappa/real_time_chat/chat-db/result"
	"github.com/WadeCappa/real_time_chat/chat-db/store"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/publisher"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	DEFAULT_PORT            = 50052
	DEFAULT_LOAD_BATCH_SIZE = 10
)

var (
	port = flag.Int("port", DEFAULT_PORT, "port for this service")
)

func getChannelManagerHostname() string {
	return os.Getenv("CHANNEL_MANAGER_HOSTNAME")
}

func getAuthHostname() string {
	return os.Getenv("AUTHMASTER_HOSTNAME")
}

func getPostgresUrl() string {
	postgresHostname := os.Getenv("CHANNEL_MANAGER_POSTGRES_HOSTNAME")
	return fmt.Sprintf("postgres://postgres:pass@%s/chat_db", postgresHostname)
}

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

func canUserWriteToChannel(channelId, userId int64) error {
	conn, err := grpc.NewClient(
		getChannelManagerHostname(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	newContext := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	c := external_channel_manager.NewExternalchannelmanagerClient(conn)
	_, err = c.CanWrite(
		newContext,
		&external_channel_manager.CanWriteRequest{
			ChannelId: channelId,
			UserId:    userId})
	if err != nil {
		return fmt.Errorf("failed permission check: %v", err)
	}

	return nil
}

func (s *chatDbServer) PublishMessage(
	ctx context.Context,
	request *chat_db.PublishMessageRequest) (*chat_db.PublishMessageResponse, error) {
	userId, err := auth.AuthenticateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %v", err)
	}

	err = canUserWriteToChannel(request.ChannelId, *userId)
	if err != nil {
		return nil, fmt.Errorf("user does not have access: %v", err)
	}

	offset, err := publisher.PublishChatMessageToChannel(*userId, request.Message, request.ChannelId)
	if err != nil {
		return nil, fmt.Errorf("failed to write to kafka: %v", err)
	}

	res := store.Call(getPostgresUrl(), func(c *pgx.Conn) result.Result[any] {
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

	result := store.Call(getPostgresUrl(), func(c *pgx.Conn) result.Result[[]chat_db.ReadMostRecentResponse] {

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
