import * as jspb from 'google-protobuf'



export class CreateChannelRequest extends jspb.Message {
  getPublic(): boolean;
  setPublic(value: boolean): CreateChannelRequest;

  getName(): string;
  setName(value: string): CreateChannelRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateChannelRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateChannelRequest): CreateChannelRequest.AsObject;
  static serializeBinaryToWriter(message: CreateChannelRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateChannelRequest;
  static deserializeBinaryFromReader(message: CreateChannelRequest, reader: jspb.BinaryReader): CreateChannelRequest;
}

export namespace CreateChannelRequest {
  export type AsObject = {
    pb_public: boolean,
    name: string,
  }
}

export class CreateChannelResponse extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): CreateChannelResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateChannelResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateChannelResponse): CreateChannelResponse.AsObject;
  static serializeBinaryToWriter(message: CreateChannelResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateChannelResponse;
  static deserializeBinaryFromReader(message: CreateChannelResponse, reader: jspb.BinaryReader): CreateChannelResponse;
}

export namespace CreateChannelResponse {
  export type AsObject = {
    channelid: number,
  }
}

export class DeleteChannelRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): DeleteChannelRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteChannelRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteChannelRequest): DeleteChannelRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteChannelRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteChannelRequest;
  static deserializeBinaryFromReader(message: DeleteChannelRequest, reader: jspb.BinaryReader): DeleteChannelRequest;
}

export namespace DeleteChannelRequest {
  export type AsObject = {
    channelid: number,
  }
}

export class DeleteChannelResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteChannelResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteChannelResponse): DeleteChannelResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteChannelResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteChannelResponse;
  static deserializeBinaryFromReader(message: DeleteChannelResponse, reader: jspb.BinaryReader): DeleteChannelResponse;
}

export namespace DeleteChannelResponse {
  export type AsObject = {
  }
}

export class ChangeChannelVisibilityRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): ChangeChannelVisibilityRequest;

  getPublic(): boolean;
  setPublic(value: boolean): ChangeChannelVisibilityRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeChannelVisibilityRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeChannelVisibilityRequest): ChangeChannelVisibilityRequest.AsObject;
  static serializeBinaryToWriter(message: ChangeChannelVisibilityRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeChannelVisibilityRequest;
  static deserializeBinaryFromReader(message: ChangeChannelVisibilityRequest, reader: jspb.BinaryReader): ChangeChannelVisibilityRequest;
}

export namespace ChangeChannelVisibilityRequest {
  export type AsObject = {
    channelid: number,
    pb_public: boolean,
  }
}

export class ChangeChannelVisibilityResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeChannelVisibilityResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeChannelVisibilityResponse): ChangeChannelVisibilityResponse.AsObject;
  static serializeBinaryToWriter(message: ChangeChannelVisibilityResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeChannelVisibilityResponse;
  static deserializeBinaryFromReader(message: ChangeChannelVisibilityResponse, reader: jspb.BinaryReader): ChangeChannelVisibilityResponse;
}

export namespace ChangeChannelVisibilityResponse {
  export type AsObject = {
  }
}

export class JoinChannelRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): JoinChannelRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinChannelRequest.AsObject;
  static toObject(includeInstance: boolean, msg: JoinChannelRequest): JoinChannelRequest.AsObject;
  static serializeBinaryToWriter(message: JoinChannelRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinChannelRequest;
  static deserializeBinaryFromReader(message: JoinChannelRequest, reader: jspb.BinaryReader): JoinChannelRequest;
}

export namespace JoinChannelRequest {
  export type AsObject = {
    channelid: number,
  }
}

export class JoinChannelResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinChannelResponse.AsObject;
  static toObject(includeInstance: boolean, msg: JoinChannelResponse): JoinChannelResponse.AsObject;
  static serializeBinaryToWriter(message: JoinChannelResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinChannelResponse;
  static deserializeBinaryFromReader(message: JoinChannelResponse, reader: jspb.BinaryReader): JoinChannelResponse;
}

export namespace JoinChannelResponse {
  export type AsObject = {
  }
}

export class AddToChannelRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): AddToChannelRequest;

  getUserid(): number;
  setUserid(value: number): AddToChannelRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddToChannelRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddToChannelRequest): AddToChannelRequest.AsObject;
  static serializeBinaryToWriter(message: AddToChannelRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddToChannelRequest;
  static deserializeBinaryFromReader(message: AddToChannelRequest, reader: jspb.BinaryReader): AddToChannelRequest;
}

export namespace AddToChannelRequest {
  export type AsObject = {
    channelid: number,
    userid: number,
  }
}

export class AddToChannelResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddToChannelResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddToChannelResponse): AddToChannelResponse.AsObject;
  static serializeBinaryToWriter(message: AddToChannelResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddToChannelResponse;
  static deserializeBinaryFromReader(message: AddToChannelResponse, reader: jspb.BinaryReader): AddToChannelResponse;
}

export namespace AddToChannelResponse {
  export type AsObject = {
  }
}

export class GetChannelsRequest extends jspb.Message {
  getPrefixsearch(): string;
  setPrefixsearch(value: string): GetChannelsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetChannelsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetChannelsRequest): GetChannelsRequest.AsObject;
  static serializeBinaryToWriter(message: GetChannelsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetChannelsRequest;
  static deserializeBinaryFromReader(message: GetChannelsRequest, reader: jspb.BinaryReader): GetChannelsRequest;
}

export namespace GetChannelsRequest {
  export type AsObject = {
    prefixsearch: string,
  }
}

export class GetChannelsResponse extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): GetChannelsResponse;

  getChannelname(): string;
  setChannelname(value: string): GetChannelsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetChannelsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetChannelsResponse): GetChannelsResponse.AsObject;
  static serializeBinaryToWriter(message: GetChannelsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetChannelsResponse;
  static deserializeBinaryFromReader(message: GetChannelsResponse, reader: jspb.BinaryReader): GetChannelsResponse;
}

export namespace GetChannelsResponse {
  export type AsObject = {
    channelid: number,
    channelname: string,
  }
}

export class GetAllChannelsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllChannelsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllChannelsRequest): GetAllChannelsRequest.AsObject;
  static serializeBinaryToWriter(message: GetAllChannelsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllChannelsRequest;
  static deserializeBinaryFromReader(message: GetAllChannelsRequest, reader: jspb.BinaryReader): GetAllChannelsRequest;
}

export namespace GetAllChannelsRequest {
  export type AsObject = {
  }
}

export class GetAllChannelsResponse extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): GetAllChannelsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllChannelsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllChannelsResponse): GetAllChannelsResponse.AsObject;
  static serializeBinaryToWriter(message: GetAllChannelsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllChannelsResponse;
  static deserializeBinaryFromReader(message: GetAllChannelsResponse, reader: jspb.BinaryReader): GetAllChannelsResponse;
}

export namespace GetAllChannelsResponse {
  export type AsObject = {
    channelid: number,
  }
}

export class CanWatchRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): CanWatchRequest;

  getUserid(): number;
  setUserid(value: number): CanWatchRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CanWatchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CanWatchRequest): CanWatchRequest.AsObject;
  static serializeBinaryToWriter(message: CanWatchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CanWatchRequest;
  static deserializeBinaryFromReader(message: CanWatchRequest, reader: jspb.BinaryReader): CanWatchRequest;
}

export namespace CanWatchRequest {
  export type AsObject = {
    channelid: number,
    userid: number,
  }
}

export class CanWatchResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CanWatchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CanWatchResponse): CanWatchResponse.AsObject;
  static serializeBinaryToWriter(message: CanWatchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CanWatchResponse;
  static deserializeBinaryFromReader(message: CanWatchResponse, reader: jspb.BinaryReader): CanWatchResponse;
}

export namespace CanWatchResponse {
  export type AsObject = {
  }
}

export class CanWriteRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): CanWriteRequest;

  getUserid(): number;
  setUserid(value: number): CanWriteRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CanWriteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CanWriteRequest): CanWriteRequest.AsObject;
  static serializeBinaryToWriter(message: CanWriteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CanWriteRequest;
  static deserializeBinaryFromReader(message: CanWriteRequest, reader: jspb.BinaryReader): CanWriteRequest;
}

export namespace CanWriteRequest {
  export type AsObject = {
    channelid: number,
    userid: number,
  }
}

export class CanWriteResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CanWriteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CanWriteResponse): CanWriteResponse.AsObject;
  static serializeBinaryToWriter(message: CanWriteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CanWriteResponse;
  static deserializeBinaryFromReader(message: CanWriteResponse, reader: jspb.BinaryReader): CanWriteResponse;
}

export namespace CanWriteResponse {
  export type AsObject = {
  }
}

