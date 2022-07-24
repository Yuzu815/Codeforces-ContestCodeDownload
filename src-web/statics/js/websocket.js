document.addEventListener("DOMContentLoaded", event => {
    const showDom = document.querySelector('#txtShow')
    const url = 'ws://localhost:8080/realtime_ws'
    const ws = new WebSocket(url);
    ws.onopen = () => {
        console.log('open connection')
    }
    ws.onclose = () => {
        console.log('close connection');
    }
    ws.onmessage = event => {
        let txt = event.data
        if (showDom.innerHTML === "NULL") {
            showDom.innerHTML = txt
        } else {
            showDom.innerHTML = showDom.innerHTML + "\n" + txt
        }
    }
});