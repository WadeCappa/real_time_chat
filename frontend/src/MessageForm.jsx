import { useState } from "react";
import { PostMessage } from "./Model";

export default function MessageForm({setter}) {
    const [message, setMessage] = useState("");

    const handleSubmit = (event) => {
        event.preventDefault();
        PostMessage(setter, message)
    }

    return (
      <form onSubmit={handleSubmit}>
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