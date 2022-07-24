function polling(){
    const eps = 1e-8
    const url = '/process';
    const xmlHttpRequest = new XMLHttpRequest();
    xmlHttpRequest.responseType = "text";
    xmlHttpRequest.open('GET', url);
    xmlHttpRequest.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;")
    xmlHttpRequest.send();
    xmlHttpRequest.onreadystatechange = function (e) {
        if (xmlHttpRequest.readyState == 4 && xmlHttpRequest.status == 200) {
            const responseText = xmlHttpRequest.responseText;
            const processBarDiv = document.getElementById("processBarInfo");
            const tempDebug = document.getElementById("tempDebug");
            const processBarDivJson = JSON.parse(responseText)
            processBarDiv.style.width = processBarDivJson.CurrentMissionProgress + "%";
            tempDebug.innerHTML = responseText
            console.log(processBarDivJson.CurrentMissionProgress )
            processBarDiv.innerHTML = processBarDivJson.CurrentMissionProgress.substring(0, 4) + "%"
            if (100.0 - parseFloat(processBarDivJson.CurrentMissionProgress) <= eps) {
                clearInterval(intervalID);
                alert("Mission accomplished. Close the window will redirect you to the code zip download page.");
                window.location.replace("/download");
            }
        }
    };
}
const intervalID = setInterval("polling()", "1000");