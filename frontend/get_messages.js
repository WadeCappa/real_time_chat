function getHello() {
    const url = 'https://api.cantseewater.online/get'
    fetch(url)
    .then(response => response.json())  
    .then(json => {
        console.log(json);
    })
}
