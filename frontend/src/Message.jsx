
export default function Message({content, timeAsMillis}) {
  return (
    <div>
        <p>posted at {timeAsMillis}</p>
        <p>{content}</p>
    </div>
  );
}
