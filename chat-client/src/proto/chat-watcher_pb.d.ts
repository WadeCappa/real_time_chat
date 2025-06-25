import * as jspb from 'google-protobuf'



export class WatchChannelResponse extends jspb.Message {
  getEvent(): ChannelEvent | undefined;
  setEvent(value?: ChannelEvent): WatchChannelResponse;
  hasEvent(): boolean;
  clearEvent(): WatchChannelResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WatchChannelResponse.AsObject;
  static toObject(includeInstance: boolean, msg: WatchChannelResponse): WatchChannelResponse.AsObject;
  static serializeBinaryToWriter(message: WatchChannelResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WatchChannelResponse;
  static deserializeBinaryFromReader(message: WatchChannelResponse, reader: jspb.BinaryReader): WatchChannelResponse;
}

export namespace WatchChannelResponse {
  export type AsObject = {
    event?: ChannelEvent.AsObject,
  }
}

export class WatchChannelRequest extends jspb.Message {
  getChannelid(): number;
  setChannelid(value: number): WatchChannelRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WatchChannelRequest.AsObject;
  static toObject(includeInstance: boolean, msg: WatchChannelRequest): WatchChannelRequest.AsObject;
  static serializeBinaryToWriter(message: WatchChannelRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WatchChannelRequest;
  static deserializeBinaryFromReader(message: WatchChannelRequest, reader: jspb.BinaryReader): WatchChannelRequest;
}

export namespace WatchChannelRequest {
  export type AsObject = {
    channelid: number,
  }
}

export class NewMessageEvent extends jspb.Message {
  getConent(): string;
  setConent(value: string): NewMessageEvent;

  getUserid(): number;
  setUserid(value: number): NewMessageEvent;

  getChannelid(): number;
  setChannelid(value: number): NewMessageEvent;

  getMessageid(): number;
  setMessageid(value: number): NewMessageEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewMessageEvent.AsObject;
  static toObject(includeInstance: boolean, msg: NewMessageEvent): NewMessageEvent.AsObject;
  static serializeBinaryToWriter(message: NewMessageEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewMessageEvent;
  static deserializeBinaryFromReader(message: NewMessageEvent, reader: jspb.BinaryReader): NewMessageEvent;
}

export namespace NewMessageEvent {
  export type AsObject = {
    conent: string,
    userid: number,
    channelid: number,
    messageid: number,
  }
}

export class UnknownEvent extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): UnknownEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnknownEvent.AsObject;
  static toObject(includeInstance: boolean, msg: UnknownEvent): UnknownEvent.AsObject;
  static serializeBinaryToWriter(message: UnknownEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnknownEvent;
  static deserializeBinaryFromReader(message: UnknownEvent, reader: jspb.BinaryReader): UnknownEvent;
}

export namespace UnknownEvent {
  export type AsObject = {
    description: string,
  }
}

export class ChannelEvent extends jspb.Message {
  getNewmessage(): NewMessageEvent | undefined;
  setNewmessage(value?: NewMessageEvent): ChannelEvent;
  hasNewmessage(): boolean;
  clearNewmessage(): ChannelEvent;

  getUnknownevent(): UnknownEvent | undefined;
  setUnknownevent(value?: UnknownEvent): ChannelEvent;
  hasUnknownevent(): boolean;
  clearUnknownevent(): ChannelEvent;

  getTimepostedunixtime(): number;
  setTimepostedunixtime(value: number): ChannelEvent;

  getOffest(): number;
  setOffest(value: number): ChannelEvent;

  getEventunionCase(): ChannelEvent.EventunionCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChannelEvent.AsObject;
  static toObject(includeInstance: boolean, msg: ChannelEvent): ChannelEvent.AsObject;
  static serializeBinaryToWriter(message: ChannelEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChannelEvent;
  static deserializeBinaryFromReader(message: ChannelEvent, reader: jspb.BinaryReader): ChannelEvent;
}

export namespace ChannelEvent {
  export type AsObject = {
    newmessage?: NewMessageEvent.AsObject,
    unknownevent?: UnknownEvent.AsObject,
    timepostedunixtime: number,
    offest: number,
  }

  export enum EventunionCase { 
    EVENTUNION_NOT_SET = 0,
    NEWMESSAGE = 1,
    UNKNOWNEVENT = 2,
  }
}

