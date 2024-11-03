import { useEffect, useState } from "react"
import getMessages from './Model';
import Message from "./Message";

export function AllMessages() {
  const [messages, setMessages] = useState([])

  useEffect(() => {
    getMessages.then((res) => setMessages(res))
  })

  return (
    {messages.map((message) => (
      <div>{message}</div>
    ))}
  )
}