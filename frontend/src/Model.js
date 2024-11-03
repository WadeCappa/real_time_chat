const apiUrl = 'https://api.cantseewater.online';

export default function getMessages() {
    return fetch(apiUrl + '/get')
        .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to get messages from backend');
                }
            return response.json();
        })
    .then(data => {
        console.log(data);
    })
    .catch(error => {
        console.error('Error:', error);
    });
}