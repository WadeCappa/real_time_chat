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
  }

  const handleBoxClicked = (event) => {
    const postId = event.target.id;
    if (event.target.checked) {
      setSelectedMessageIds([...selectedMessageIds, postId])
    } else {
      setSelectedMessageIds(selectedMessageIds.filter(p => p !== postId))
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
                    <input type="checkbox" id={message.PostId} name={message.PostId} onChange={handleBoxClicked}/>
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