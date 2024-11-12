import { useState } from "react";
import { PostMessage } from "./Model";

export default function MessageForm() {
    const [message, setMessage] = useState("");

    const handleSubmit = (event) => {
        event.preventDefault();
        PostMessage(message, setMessage)
    }

    return (
      <form onSubmit={handleSubmit} style={{margin: '8px', marginBottom: '20px'}}>
        <label> Enter your message
          <input 
            type="text" 
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
        </label>
        <input type="submit" />
      </form>
    )
  }