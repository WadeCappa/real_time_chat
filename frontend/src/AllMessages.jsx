import { useEffect, useState } from "react"
import { GetMessages } from "./Model";
import MessageForm from "./MessageForm";

export default function AllMessages() {
  const [messages, setMessages] = useState([])

  useEffect(() => GetMessages(setMessages), [])

  return (
    <div>
      <MessageForm setter={setMessages} />
      <table >
        <tr>
          <th>Message</th>
          <th>Time Posted</th>
        </tr>
        {messages ? messages.map(message => {
          return (
            <tr>
              <td>{message.Content}</td>
              <td>{message.TimePosted}</td>
            </tr>
          )
        }) : 'Loading...'}
      </table> 
    </div>
  )
}