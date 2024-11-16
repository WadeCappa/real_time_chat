import { useRef, useEffect, useState } from "react"
import { GetMessages, DeleteMessages, WatchForNewMessages } from "./Model";
import MessageForm from "./MessageForm";

export default function AllMessages() {
  const [messages, setMessages] = useState([])
  const [selectedMessageIds, setSelectedMessageIds] = useState([])
  const messagesRef = useRef(messages)

  useEffect(() => {
    GetMessages((m) => {
      messagesRef.current = m
      setMessages(m)
    });
    WatchForNewMessages((m) => {
      console.log(m)
      console.log(messagesRef.current)
      messagesRef.current = [m, ...messagesRef.current]
      setMessages(messagesRef.current)
    }, (e) => {
      console.log(e)
      console.log(messagesRef.current)
      messagesRef.current = messagesRef.current.filter(m => m.PostId !== e)
      setMessages(messagesRef.current)
    })
  }, [])

  const handleSubmit = (event) => {
    event.preventDefault();
    DeleteMessages(selectedMessageIds)
    setSelectedMessageIds([])
  }

  const handleBoxClicked = (event) => {
    const postId = event.target.id;
    if (selectedMessageIds.includes(postId)) {
      setSelectedMessageIds(selectedMessageIds.filter(p => p !== postId))
    } else {
      setSelectedMessageIds([...selectedMessageIds, postId])
    }
    console.log(selectedMessageIds)
  }

  return (
    <div>
      <MessageForm />
      <form onSubmit={handleSubmit} id="choosePostForm">
        <button type="submit" form="choosePostForm" value="Submit">Delete selected</button>
        <table >
          <tr>
            <th>Selected</th>
            <th>Message</th>
            <th>Time Posted</th>
          </tr>
          {messages ? messages.map(message => {
            return (
                <tr>
                  <td>
                    <input 
                      type="checkbox" 
                      id={message.PostId}
                      onChange={handleBoxClicked}
                    />
                  </td>
                  <td>
                    <label for={message.PostId}>
                      {message.Content}
                    </label>
                  </td>
                  <td>
                    <label for={message.PostId}>
                      {message.TimePosted}
                    </label>
                  </td>
                </tr>
            )
          }) : "Loading..."}
        </table> 
      </form>
    </div>
  )
}