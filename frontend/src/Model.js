
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

export function GetMessages(setter) {
    const url = getUrl()
    fetch(url + '/get')
    .then(response => {
        console.log(response)
        return response.json()
    })
    .then(data => setter(data))
    .catch(error => {
        console.error('Error:', error);
    });
}

export function PostMessage(setter, newMessage) {
    const url = getUrl()
    const request = {
        method: "POST",
        body: JSON.stringify({"content": newMessage}),
    }
    fetch(url + '/write', request)
    .then(response => {
        console.log(response)
        return response.json()
    })
    .then(_ => GetMessages(setter))
    .catch(error => {
        console.error('Error:', error);
    });
}