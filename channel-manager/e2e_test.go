package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/WadeCappa/authmaster/authmaster"
	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	serviceHostname     = fmt.Sprintf("localhost:%d", DEFAULT_PORT)
	testingAuthHostname = "localhost:50051"
)

func createRandomString() string {
	return fmt.Sprintf("%d", rand.Int())
}

type user struct {
	token  string
	userId int64
}

func makeFakeUser() user {
	user := user{}

	password := createRandomString()
	username := createRandomString()

	err := withConnection(testingAuthHostname, func(cc *grpc.ClientConn) error {
		c := authmaster.NewAuthmasterClient(cc)
		_, err := c.CreateUser(
			context.Background(),
			&authmaster.CreateUserRequest{
				Username: username,
				Password: password})

		if err != nil {
			log.Fatalf("failed to create user %v", err)
		}

		login, err := c.Login(
			context.Background(),
			&authmaster.LoginRequest{
				Username: username,
				Password: password})
		if err != nil {
			log.Fatalf("failed to login %v", err)
		}

		user.token = login.Token

		res, err := c.TestAuth(
			metadata.NewOutgoingContext(
				context.Background(),
				metadata.Pairs("Authorization", user.token)),
			&authmaster.TestAuthRequest{})

		if err != nil {
			log.Fatalf("failed to test auth %v", err)
		}

		user.userId = res.UserId
		return nil
	})

	if err != nil {
		log.Fatalf("failed! %v", err)
	}
	return user
}

func withConnection(addr string, consumer func(*grpc.ClientConn) error) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	return consumer(conn)
}

func checkWritePermissions(channelId, userId int64, c external_channel_manager.ExternalchannelmanagerClient) error {
	_, err := c.CanWrite(
		context.Background(),
		&external_channel_manager.CanWriteRequest{
			UserId:    userId,
			ChannelId: channelId})
	return err
}

func checkReadPermissions(channelId, userId int64, c external_channel_manager.ExternalchannelmanagerClient) error {
	_, err := c.CanWatch(
		context.Background(),
		&external_channel_manager.CanWatchRequest{
			UserId:    userId,
			ChannelId: channelId})
	return err
}

func TestMakingChannel(t *testing.T) {

	user := makeFakeUser()

	withConnection(serviceHostname, func(cc *grpc.ClientConn) error {
		newMetadata := metadata.Pairs("Authorization", user.token)
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := external_channel_manager.NewExternalchannelmanagerClient(cc)

		_, err := c.CreateChannel(
			newContext,
			&external_channel_manager.CreateChannelRequest{
				Name:   createRandomString(),
				Public: true,
			})
		if err != nil {
			t.Errorf("failed to create channel %v", err)
		}

		return nil
	})
}

func TestJoinPublicChannel(t *testing.T) {
	author := makeFakeUser()

	withConnection(serviceHostname, func(cc *grpc.ClientConn) error {
		authorContext := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("Authorization", author.token))

		c := external_channel_manager.NewExternalchannelmanagerClient(cc)

		resp, err := c.CreateChannel(
			authorContext,
			&external_channel_manager.CreateChannelRequest{
				Name:   createRandomString(),
				Public: true,
			})
		if err != nil {
			t.Errorf("failed to create channel %v", err)
		}

		channelId := resp.ChannelId

		otherUser := makeFakeUser()

		if checkWritePermissions(channelId, otherUser.userId, c) == nil {
			t.Error("should not have write permissions")
		}

		if checkReadPermissions(channelId, otherUser.userId, c) != nil {
			t.Error("should have read permissions")
		}

		otherUserContext := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("Authorization", otherUser.token))

		_, err = c.JoinChannel(
			otherUserContext,
			&external_channel_manager.JoinChannelRequest{
				ChannelId: channelId})
		if err != nil {
			t.Errorf("failed to join channel %v", err)
		}

		if checkWritePermissions(channelId, otherUser.userId, c) != nil {
			t.Error("should have write permissions")
		}

		if checkReadPermissions(channelId, otherUser.userId, c) != nil {
			t.Error("should still have read permissions")
		}
		return nil
	})
}

func TestCantJoinPrivateChannel(t *testing.T) {
	author := makeFakeUser()

	withConnection(serviceHostname, func(cc *grpc.ClientConn) error {
		authorContext := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("Authorization", author.token))

		c := external_channel_manager.NewExternalchannelmanagerClient(cc)

		resp, err := c.CreateChannel(
			authorContext,
			&external_channel_manager.CreateChannelRequest{
				Name:   createRandomString(),
				Public: false,
			})
		if err != nil {
			t.Errorf("failed to create channel %v", err)
		}

		channelId := resp.ChannelId

		otherUser := makeFakeUser()

		if checkWritePermissions(channelId, otherUser.userId, c) == nil {
			t.Error("should not have write permissions")
		}

		if checkReadPermissions(channelId, otherUser.userId, c) == nil {
			t.Error("should not have read permissions")
		}

		otherUserContext := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("Authorization", otherUser.token))

		_, err = c.JoinChannel(
			otherUserContext,
			&external_channel_manager.JoinChannelRequest{
				ChannelId: channelId})

		if err == nil {
			t.Errorf("should not have been able to join channel %v", err)
		}

		if checkWritePermissions(channelId, otherUser.userId, c) == nil {
			t.Error("should still not have write permissions")
		}

		if checkReadPermissions(channelId, otherUser.userId, c) == nil {
			t.Error("should still not have read permissions")
		}

		return nil
	})
}

func TestAddUserToChannel(t *testing.T) {
	author := makeFakeUser()

	withConnection(serviceHostname, func(cc *grpc.ClientConn) error {
		authorContext := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("Authorization", author.token))

		c := external_channel_manager.NewExternalchannelmanagerClient(cc)

		resp, err := c.CreateChannel(
			authorContext,
			&external_channel_manager.CreateChannelRequest{
				Name:   createRandomString(),
				Public: false,
			})
		if err != nil {
			t.Errorf("failed to create channel %v", err)
		}

		channelId := resp.ChannelId

		otherUser := makeFakeUser()

		_, err = c.AddToChannel(
			authorContext,
			&external_channel_manager.AddToChannelRequest{
				ChannelId: channelId,
				UserId:    otherUser.userId})
		if err != nil {
			t.Errorf("should have succeeded %v", err)
		}

		if checkWritePermissions(channelId, otherUser.userId, c) != nil {
			t.Error("should have write permissions")
		}

		if checkReadPermissions(channelId, otherUser.userId, c) != nil {
			t.Error("should have read permissions")
		}

		return nil
	})
}
