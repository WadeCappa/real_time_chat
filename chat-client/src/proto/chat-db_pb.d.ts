import * as jspb from 'google-protobuf'



export class ReadMostRecentResponse extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): ReadMostRecentResponse;

  getMessage(): string;
  setMessage(value: string): ReadMostRecentResponse;

  getUserid(): number;
  setUserid(value: number): ReadMostRecentResponse;

  getTimepostedunixtime(): number;
  setTimepostedunixtime(value: number): ReadMostRecentResponse;

  getMessageid(): number;
  setMessageid(value: number): ReadMostRecentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadMostRecentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReadMostRecentResponse): ReadMostRecentResponse.AsObject;
  static serializeBinaryToWriter(message: ReadMostRecentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadMostRecentResponse;
  static deserializeBinaryFromReader(message: ReadMostRecentResponse, reader: jspb.BinaryReader): ReadMostRecentResponse;
}

export namespace ReadMostRecentResponse {
  export type AsObject = {
    channelid: number,
    message: string,
    userid: number,
    timepostedunixtime: number,
    messageid: number,
  }
}

export class ReadMostRecentRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): ReadMostRecentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadMostRecentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReadMostRecentRequest): ReadMostRecentRequest.AsObject;
  static serializeBinaryToWriter(message: ReadMostRecentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadMostRecentRequest;
  static deserializeBinaryFromReader(message: ReadMostRecentRequest, reader: jspb.BinaryReader): ReadMostRecentRequest;
}

export namespace ReadMostRecentRequest {
  export type AsObject = {
    channelid: number,
  }
}

export class PublishMessageResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishMessageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PublishMessageResponse): PublishMessageResponse.AsObject;
  static serializeBinaryToWriter(message: PublishMessageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishMessageResponse;
  static deserializeBinaryFromReader(message: PublishMessageResponse, reader: jspb.BinaryReader): PublishMessageResponse;
}

export namespace PublishMessageResponse {
  export type AsObject = {
  }
}

export class PublishMessageRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): PublishMessageRequest;

  getMessage(): string;
  setMessage(value: string): PublishMessageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishMessageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PublishMessageRequest): PublishMessageRequest.AsObject;
  static serializeBinaryToWriter(message: PublishMessageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishMessageRequest;
  static deserializeBinaryFromReader(message: PublishMessageRequest, reader: jspb.BinaryReader): PublishMessageRequest;
}

export namespace PublishMessageRequest {
  export type AsObject = {
    channelid: number,
    message: string,
  }
}

