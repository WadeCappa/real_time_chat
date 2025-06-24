package main

import (
	"context"
	"fmt"

	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"github.com/WadeCappa/real_time_chat/channel-manager/external_endpoints"
	"github.com/WadeCappa/real_time_chat/channel-manager/internal_endpoints"
	"google.golang.org/grpc"
)

type externalChatManangerServer struct {
	external_channel_manager.ExternalchannelmanagerServer
}

func (s *externalChatManangerServer) AddToChannel(ctx context.Context, request *external_channel_manager.AddToChannelRequest) (*external_channel_manager.AddToChannelResponse, error) {
	userId, err := auth.AuthenticateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return external_endpoints.AddToChannel(getPostgresUrl(), *userId, ctx, request)
}

func (s *externalChatManangerServer) ChangeChannelVisibility(ctx context.Context, request *external_channel_manager.ChangeChannelVisibilityRequest) (*external_channel_manager.ChangeChannelVisibilityResponse, error) {
	userId, err := auth.AuthenticateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return external_endpoints.ChangeChannelVisibility(getPostgresUrl(), *userId, ctx, request)
}

func (s *externalChatManangerServer) CreateChannel(ctx context.Context, request *external_channel_manager.CreateChannelRequest) (*external_channel_manager.CreateChannelResponse, error) {
	userId, err := auth.AuthenticateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return external_endpoints.CreateChannel(getPostgresUrl(), *userId, ctx, request)
}

func (s *externalChatManangerServer) DeleteChannel(ctx context.Context, request *external_channel_manager.DeleteChannelRequest) (*external_channel_manager.DeleteChannelResponse, error) {
	userId, err := auth.AuthenticateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return external_endpoints.DeleteChannel(getPostgresUrl(), *userId, ctx, request)
}

func (s *externalChatManangerServer) JoinChannel(ctx context.Context, request *external_channel_manager.JoinChannelRequest) (*external_channel_manager.JoinChannelResponse, error) {
	userId, err := auth.AuthenticateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return external_endpoints.JoinChannel(getPostgresUrl(), *userId, ctx, request)
}

func (s *externalChatManangerServer) GetChannels(request *external_channel_manager.GetChannelsRequest, server grpc.ServerStreamingServer[external_channel_manager.GetChannelsResponse]) error {
	userId, err := auth.AuthenticateUser(server.Context())
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	return external_endpoints.GetChannels(getPostgresUrl(), *userId, request, server)
}

func (s *externalChatManangerServer) CanWrite(ctx context.Context, request *external_channel_manager.CanWriteRequest) (*external_channel_manager.CanWriteResponse, error) {
	return internal_endpoints.CanWrite(getPostgresUrl(), ctx, request)
}

func (s *externalChatManangerServer) CanWatch(ctx context.Context, request *external_channel_manager.CanWatchRequest) (*external_channel_manager.CanWatchResponse, error) {
	return internal_endpoints.CanWatch(getPostgresUrl(), ctx, request)
}

func (s *externalChatManangerServer) GetAllChannels(request *external_channel_manager.GetAllChannelsRequest, server grpc.ServerStreamingServer[external_channel_manager.GetAllChannelsResponse]) error {
	return internal_endpoints.GetAllChannels(getPostgresUrl(), request, server)
}
