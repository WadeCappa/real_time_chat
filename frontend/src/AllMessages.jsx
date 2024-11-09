import { useEffect, useState } from "react"
import { GetMessages, DeleteMessages } from "./Model";
import MessageForm from "./MessageForm";

export default function AllMessages() {
  const [messages, setMessages] = useState([])
  const [selectedMessageIds, setSelectedMessageIds] = useState([])

  useEffect(() => GetMessages(setMessages), [])

  const handleSubmit = (event) => {
    event.preventDefault();
    DeleteMessages(setMessages, selectedMessageIds)
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
      <MessageForm setter={setMessages}  />
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
                      checked={selectedMessageIds.includes(message.PostId)} 
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