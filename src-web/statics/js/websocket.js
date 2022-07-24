document.addEventListener("DOMContentLoaded", event => {
    const url = '/realtime_ws'
    const ws = new WebSocket(wsURL(url));
    ws.onopen = () => {
        console.log('open connection')
    }
    ws.onclose = () => {
        console.log('close connection');
    }
    ws.onmessage = event => {
        const result = event.data
        const logViewer = document.getElementById('logViewer')
        var logViewerDiv = document.getElementById("logViewerDiv")
        logViewer.innerHTML = logViewer.innerHTML + result
        logViewerDiv.scrollTop=logViewerDiv.scrollHeight
    }
});


//Copy from https://www.itranslater.com/qa/details/2123328733623878656
function wsURL(s) {
    const l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + l.pathname + s;
}
