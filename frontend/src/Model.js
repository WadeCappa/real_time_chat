
const apiUrl = 'https://api.cantseewater.online';
const devUrl = 'http://localhost:8080'

function getUrl() {
    switch (process.env.REACT_APP_DEPLOYMENT) {
        case "production":
            return apiUrl;
        case "local":
            return devUrl;
        default:
            console.error(process.env.REACT_APP_DEPLOYMENT)
            throw new Error("Unrecognized deployement type")
    }
}

export function WatchForNewMessages(singleMessageSetter) {
    const eventSource = new EventSource(getUrl() + "/watch-messages")
    eventSource.onmessage = (event) => {
        const data = JSON.parse(event.data)
        singleMessageSetter(data)
    }
    return () => eventSource.close();
}

export function GetMessages(allDataSetter) {
    const url = getUrl()
    fetch(url + '/')
    .then(response => {
        console.log(response)
        return response.json()
    })
    .then(data => allDataSetter(data))
    .catch(error => {
        console.error('Error:', error);
    });
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