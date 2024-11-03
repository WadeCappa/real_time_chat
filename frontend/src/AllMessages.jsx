import { useEffect, useState } from "react"
import Message from "./Message";
import { GetMessages } from "./Model";
import MessageForm from "./MessageForm";

export default function AllMessages() {
  const [messages, setMessages] = useState([])

  useEffect(() => GetMessages(setMessages), [])

  return (
    <div>
      <MessageForm setter={setMessages}/>
      {messages ? <pre>{JSON.stringify(messages, null, 2)}</pre> : 'Loading...'}
    </div>
  )
}