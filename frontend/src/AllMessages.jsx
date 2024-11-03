import { useEffect, useState } from "react"
import Message from "./Message";

const apiUrl = 'https://api.cantseewater.online';
const devUrl = 'http://localhost:8080'

export default function AllMessages() {
  const [messages, setMessages] = useState([])

  useEffect(() => {
    const url = process.env.NODE_ENV === "production" ? apiUrl : devUrl
    fetch(url + '/get')
        .then(response => {
            console.log(response)
            return response.json()
        })
        .then(data => setMessages(data))
    .catch(error => {
        console.error('Error:', error);
    });
  }, [])

  return (
    <div>
      {messages ? <pre>{JSON.stringify(messages, null, 2)}</pre> : 'Loading...'}
    </div>
  )
}