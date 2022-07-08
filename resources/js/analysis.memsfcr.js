let serverUrl = "http://127.0.0.1:8081/scenario/contents/"
let analysisUrl = "http://127.0.0.1:5000"
let filename = "fullrun.fcr"

getFiledata()
    .then(response => {
        sendFile(response, filename)
    })


async function sendFile(data, filename) {
    let url = analysisUrl
    let form = new FormData()
    let filetype = "text/csv"

    if (filename.toLowerCase().includes(".fcr")) {
        filetype = "text/json"
    }

    const blob = new Blob([data]);
    const f = new File([blob], filename, {type: filetype});
    form.append('file', f)

    const oReq = new XMLHttpRequest();
    oReq.open("POST", url, true);
    oReq.setRequestHeader("Cache-Control", "no-cache");
    oReq.setRequestHeader("X-Requested-With", "XMLHttpRequest");

    oReq.onload = function(oEvent) {
        if (oReq.status == 200) {
            console.info("uploaded " + filename + " (" + oReq.response + ")")
        } else {
            console.error("Error " + oReq.status + " occurred");
        }
    };

    oReq.send(form);
}

async function getFiledata() {
    let url = serverUrl + filename
    let init = {
        method: 'GET',
    }

    try {
        let response = await fetch(url, init);
        return await response.text();
    } catch (e) {
        alert(e)
    }
}
