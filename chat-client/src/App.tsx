import React, { useState } from "react";
import { PublishMessageRequest } from "./proto/chat-db_pb"; 
import { GetChannelsRequest } from "./proto/external-channel-manager_pb"; 
import { externalchannelmanagerClient } from "./proto/External-channel-managerServiceClientPb";
import { chatdbClient } from "./proto/Chat-dbServiceClientPb"; 
import { chatwatcherserverClient } from "./proto/Chat-watcherServiceClientPb"; 
import { WatchChannelRequest } from "./proto/chat-watcher_pb";
import { Message } from "./Message";

const EnvoyURL = "http://localhost:8000"

const publishMessage = async (name: string, token: string, channelId: number) => {
  const client = new chatdbClient(EnvoyURL);
  const request = new PublishMessageRequest();
  request.setMessage(name);
  request.setChannelid(channelId);
  const metadata = {'authorization': token};
  const response = await client.publishMessage(request, metadata);
  console.log(response);
  const div = document.getElementById("response");
  if (div) div.innerText = JSON.stringify(response);
};

const isNumber = (v: string): boolean => {
  return !isNaN(Number(v))
}

function App() {
  const [name, setName] = useState("");
  const [token, setToken] = useState("");
  const [selectedChannelId, setSelectedChannelId] = useState("");
  const [channels, setChannels] = useState<number[]>([]);
  const [messages, setMessages] = useState<Message[]>([])

  const findChannels = async (token: string) => {
    const client = new externalchannelmanagerClient(EnvoyURL);
    const request = new GetChannelsRequest();
    const metadata = {'authorization': token};
    const response = await client.getChannels(request, metadata);
    response.on('data', resp => {
      console.log(resp.getChannelid())
      setChannels(prevChannels => [...prevChannels, resp.getChannelid()])
    })
    response.on('error', err => {
      console.log(err)
    })
  }

  const watchChat = async (token: string, channelId: number) => {
    const client = new chatwatcherserverClient(EnvoyURL);
    const request = new WatchChannelRequest();
    request.setChannelid(channelId)
    const metadata = {'authorization': token};
    const response = await client.watchChannel(request, metadata);
    console.log("before trying to listen for channel " + channelId)
    response.on('status', status => {
      console.log(JSON.stringify(status))
    })
    response.on('metadata', meta => {
      console.log(JSON.stringify(meta))
    })
    response.on('end', () => {
      // do nothing
    })
    response.on('data', resp => {
      console.log(JSON.stringify(resp))
      if (resp.getEvent() !== undefined) {
        if (resp.getEvent()!.getNewmessage() !== undefined) {
          const content = resp.getEvent()!.getNewmessage()!.getConent()
          const userId = resp.getEvent()!.getNewmessage()!.getUserid()
          const posted = resp.getEvent()!.getTimepostedunixtime()
          console.log(resp.getEvent()!.getTimepostedunixtime())
          const newMessage: Message = {userId: userId, content: content, posted: new Date(posted)}
          setMessages(prevMessages => [...prevMessages, newMessage])
        }
      }
    })
    response.on('error', err => {
      console.log(err)
    })
  }

  const onTokenUpdateFindChannels = (token: string) => {
    setChannels([])
    setToken(token)
    if (token) findChannels(token);
  }
  const onClickGreet = () => {
    if (isNumber(selectedChannelId) && name) {
      publishMessage(name, token, Number(selectedChannelId))
    }
  };
  const updateTargetChannel = (channel: string) => {
    setMessages([])
    setSelectedChannelId(channel)
    if (isNumber(selectedChannelId) && token) {
      watchChat(token, Number(channel))
    }
  }

  return (
    <div className="App">
      <div>
        Token: <input
          type="text"
          value={token}
          onChange={(e) => onTokenUpdateFindChannels(e.target.value)}
        />
      </div>
      <div>
        {channels.join(", ")}
      </div>
      <div>
        channel id: <input
          type="text"
          value={selectedChannelId}
          onChange={(e) => updateTargetChannel(e.target.value)}
        />
      </div>
      <div>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <button onClick={onClickGreet}>post</button>
        {name && <div id="response"></div>}
      </div>
      <div>
        <table >
          <tr>
            <th>Time Posted</th>
            <th>User ID</th>
            <th>Message</th>
          </tr>
          {messages ? messages.map(message => {
            return (
              <tr>
                <td>
                  {message.posted.toDateString()}
                </td>
                <td>
                  {message.userId}
                </td>
                <td>
                  {message.content}
                </td>
              </tr>
            )
          }) : "Find you user token and select a channel"}
        </table>
      </div>
    </div>
  );
}

export default App;
