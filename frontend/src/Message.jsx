
export default function Message({content, time}) {
  return (
    <div>
        <p>{content}</p>
        <p>posted at {time}</p>
    </div>
  );
}
