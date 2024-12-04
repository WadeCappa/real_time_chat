
function getUrl() {
    return process.env.REACT_APP_BACKEND_URL;
}

export function WatchForNewMessages(singleMessageSetter, messageDeleter) {
    const eventSource = new EventSource(getUrl() + "/watch")
    eventSource.onmessage = (event) => {
        console.log(event)
        const data = JSON.parse(event.data)
        console.log(data)
        switch (data.Name) {
            case "newMessage":
                singleMessageSetter(data.Payload)
                break
            case "deleteMessage":
                messageDeleter(data.Payload)
                break
            default:
                console.log("unrecognized event")
        }
    }
    return () => eventSource.close();
}

export function PostMessage(newMessage, messageSetter) {
    const url = getUrl()
    const request = {
        method: "POST",
        body: JSON.stringify({"Content": newMessage}),
    }
    fetch(url + '/', request)
    .then(_ => messageSetter(""))
    .catch(error => {
        console.error('Error:', error);
    });
}

export function DeleteMessages(messageIdsToDelete) {
    console.log(messageIdsToDelete)
    const url = getUrl()
    const request = {
        method: "DELETE",
        body: JSON.stringify({"postIds": messageIdsToDelete.map(m => Number(m))}),
    }
    fetch(url + '/', request)
    .then(_ => {})
    .catch(error => {
        console.error('Error:', error);
    });
}