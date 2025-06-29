/**
 * @fileoverview gRPC-Web generated client stub for external_channel_manager
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v3.19.6
// source: external-channel-manager.proto


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as external$channel$manager_pb from './external-channel-manager_pb'; // proto import: "external-channel-manager.proto"


export class externalchannelmanagerClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'binary';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname.replace(/\/+$/, '');
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodDescriptorCreateChannel = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/CreateChannel',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.CreateChannelRequest,
    external$channel$manager_pb.CreateChannelResponse,
    (request: external$channel$manager_pb.CreateChannelRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.CreateChannelResponse.deserializeBinary
  );

  createChannel(
    request: external$channel$manager_pb.CreateChannelRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.CreateChannelResponse>;

  createChannel(
    request: external$channel$manager_pb.CreateChannelRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.CreateChannelResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.CreateChannelResponse>;

  createChannel(
    request: external$channel$manager_pb.CreateChannelRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.CreateChannelResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/CreateChannel',
        request,
        metadata || {},
        this.methodDescriptorCreateChannel,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/CreateChannel',
    request,
    metadata || {},
    this.methodDescriptorCreateChannel);
  }

  methodDescriptorDeleteChannel = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/DeleteChannel',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.DeleteChannelRequest,
    external$channel$manager_pb.DeleteChannelResponse,
    (request: external$channel$manager_pb.DeleteChannelRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.DeleteChannelResponse.deserializeBinary
  );

  deleteChannel(
    request: external$channel$manager_pb.DeleteChannelRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.DeleteChannelResponse>;

  deleteChannel(
    request: external$channel$manager_pb.DeleteChannelRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.DeleteChannelResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.DeleteChannelResponse>;

  deleteChannel(
    request: external$channel$manager_pb.DeleteChannelRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.DeleteChannelResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/DeleteChannel',
        request,
        metadata || {},
        this.methodDescriptorDeleteChannel,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/DeleteChannel',
    request,
    metadata || {},
    this.methodDescriptorDeleteChannel);
  }

  methodDescriptorChangeChannelVisibility = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/ChangeChannelVisibility',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.ChangeChannelVisibilityRequest,
    external$channel$manager_pb.ChangeChannelVisibilityResponse,
    (request: external$channel$manager_pb.ChangeChannelVisibilityRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.ChangeChannelVisibilityResponse.deserializeBinary
  );

  changeChannelVisibility(
    request: external$channel$manager_pb.ChangeChannelVisibilityRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.ChangeChannelVisibilityResponse>;

  changeChannelVisibility(
    request: external$channel$manager_pb.ChangeChannelVisibilityRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.ChangeChannelVisibilityResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.ChangeChannelVisibilityResponse>;

  changeChannelVisibility(
    request: external$channel$manager_pb.ChangeChannelVisibilityRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.ChangeChannelVisibilityResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/ChangeChannelVisibility',
        request,
        metadata || {},
        this.methodDescriptorChangeChannelVisibility,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/ChangeChannelVisibility',
    request,
    metadata || {},
    this.methodDescriptorChangeChannelVisibility);
  }

  methodDescriptorJoinChannel = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/JoinChannel',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.JoinChannelRequest,
    external$channel$manager_pb.JoinChannelResponse,
    (request: external$channel$manager_pb.JoinChannelRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.JoinChannelResponse.deserializeBinary
  );

  joinChannel(
    request: external$channel$manager_pb.JoinChannelRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.JoinChannelResponse>;

  joinChannel(
    request: external$channel$manager_pb.JoinChannelRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.JoinChannelResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.JoinChannelResponse>;

  joinChannel(
    request: external$channel$manager_pb.JoinChannelRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.JoinChannelResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/JoinChannel',
        request,
        metadata || {},
        this.methodDescriptorJoinChannel,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/JoinChannel',
    request,
    metadata || {},
    this.methodDescriptorJoinChannel);
  }

  methodDescriptorAddToChannel = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/AddToChannel',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.AddToChannelRequest,
    external$channel$manager_pb.AddToChannelResponse,
    (request: external$channel$manager_pb.AddToChannelRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.AddToChannelResponse.deserializeBinary
  );

  addToChannel(
    request: external$channel$manager_pb.AddToChannelRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.AddToChannelResponse>;

  addToChannel(
    request: external$channel$manager_pb.AddToChannelRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.AddToChannelResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.AddToChannelResponse>;

  addToChannel(
    request: external$channel$manager_pb.AddToChannelRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.AddToChannelResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/AddToChannel',
        request,
        metadata || {},
        this.methodDescriptorAddToChannel,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/AddToChannel',
    request,
    metadata || {},
    this.methodDescriptorAddToChannel);
  }

  methodDescriptorGetChannels = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/GetChannels',
    grpcWeb.MethodType.SERVER_STREAMING,
    external$channel$manager_pb.GetChannelsRequest,
    external$channel$manager_pb.GetChannelsResponse,
    (request: external$channel$manager_pb.GetChannelsRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.GetChannelsResponse.deserializeBinary
  );

  getChannels(
    request: external$channel$manager_pb.GetChannelsRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<external$channel$manager_pb.GetChannelsResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/external_channel_manager.externalchannelmanager/GetChannels',
      request,
      metadata || {},
      this.methodDescriptorGetChannels);
  }

  methodDescriptorGetAllChannels = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/GetAllChannels',
    grpcWeb.MethodType.SERVER_STREAMING,
    external$channel$manager_pb.GetAllChannelsRequest,
    external$channel$manager_pb.GetAllChannelsResponse,
    (request: external$channel$manager_pb.GetAllChannelsRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.GetAllChannelsResponse.deserializeBinary
  );

  getAllChannels(
    request: external$channel$manager_pb.GetAllChannelsRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<external$channel$manager_pb.GetAllChannelsResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/external_channel_manager.externalchannelmanager/GetAllChannels',
      request,
      metadata || {},
      this.methodDescriptorGetAllChannels);
  }

  methodDescriptorCanWatch = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/CanWatch',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.CanWatchRequest,
    external$channel$manager_pb.CanWatchResponse,
    (request: external$channel$manager_pb.CanWatchRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.CanWatchResponse.deserializeBinary
  );

  canWatch(
    request: external$channel$manager_pb.CanWatchRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.CanWatchResponse>;

  canWatch(
    request: external$channel$manager_pb.CanWatchRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.CanWatchResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.CanWatchResponse>;

  canWatch(
    request: external$channel$manager_pb.CanWatchRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.CanWatchResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/CanWatch',
        request,
        metadata || {},
        this.methodDescriptorCanWatch,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/CanWatch',
    request,
    metadata || {},
    this.methodDescriptorCanWatch);
  }

  methodDescriptorCanWrite = new grpcWeb.MethodDescriptor(
    '/external_channel_manager.externalchannelmanager/CanWrite',
    grpcWeb.MethodType.UNARY,
    external$channel$manager_pb.CanWriteRequest,
    external$channel$manager_pb.CanWriteResponse,
    (request: external$channel$manager_pb.CanWriteRequest) => {
      return request.serializeBinary();
    },
    external$channel$manager_pb.CanWriteResponse.deserializeBinary
  );

  canWrite(
    request: external$channel$manager_pb.CanWriteRequest,
    metadata?: grpcWeb.Metadata | null): Promise<external$channel$manager_pb.CanWriteResponse>;

  canWrite(
    request: external$channel$manager_pb.CanWriteRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.CanWriteResponse) => void): grpcWeb.ClientReadableStream<external$channel$manager_pb.CanWriteResponse>;

  canWrite(
    request: external$channel$manager_pb.CanWriteRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: external$channel$manager_pb.CanWriteResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/external_channel_manager.externalchannelmanager/CanWrite',
        request,
        metadata || {},
        this.methodDescriptorCanWrite,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/external_channel_manager.externalchannelmanager/CanWrite',
    request,
    metadata || {},
    this.methodDescriptorCanWrite);
  }

}

