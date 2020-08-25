var conn
function changeDifficultyLevel(){
    resetTimer()
    var level = document.getElementById("level").value
    var table = document.getElementById("grid")
    var ol = document.getElementById("scoreBoard")
    if(table != null) {
        while(table.firstChild) {
            table.removeChild(table.firstChild);
        }
    }
    while(ol.firstChild) {
        ol.removeChild(ol.firstChild);
    }
    console.log("Level selected : ",level)
    loadGame(level)
}
async function loadGame(level) {
    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.close = function () {
        console.log("Connection with Server closed")
    }
    await new Promise(r => setTimeout(r, 500));
    conn.send(level)
    conn.onmessage = function (server_data) {
        var ol = document.getElementById("scoreBoard")
        var message = server_data.data
        
        console.log("Message : ",message)
        if( message.length > 81) {
            var obj = JSON.parse(message)
            for (var i = 0; i < obj.length; i++) {
                var score = obj[i];
                var li = document.createElement("li")
                li.innerHTML = score.Name + " - " + score.Time
                ol.appendChild(li)
            }
        } else {
            if(message == null) {
                console.log("in null")
                var li = document.createElement("li")
                li.innerHTML = "No Data Present"
                ol.appendChild(li)
            }
            var tr;
            var tbl = document.getElementById("grid")
            for (var i = 0; i < message.length; i++) {
                if (i % 9 == 0) {
                    tr = document.createElement("tr")
                    tbl.appendChild(tr)
                }
                var td = document.createElement('td')
                // td.innerHTML = message[i]
                if (message[i] == '0') {
                    td.innerHTML = "<input id='cell-" + i + "'  type='number' min='1' max='9'  value='' onchange='sendMessage(" + i + ")'>"
                } else {
                    td.innerHTML = "<input id='cell-" + i + "'  type='text' value='" + message[i] + "'disabled>"
                }

                tr.appendChild(td)
            }
            timer("start")
        }
        // console.log("Server : ", message)
    }
    
}
        
window.onload = function () {
    if (window["WebSocket"]) {
        loadGame(0)
    } else {
        console.log("Your Browser does not websocket")
        alert("Your Browser does not websocket")
    }
}

        
function resetGame() {
    resetTimer()
    var table = document.getElementById("grid")
    var ol = document.getElementById("scoreBoard")
    if(table != null) {
        while(table.firstChild) {
            table.removeChild(table.firstChild);
        }
    }
    while(ol.firstChild) {
        ol.removeChild(ol.firstChild);
    }
    loadGame(0)
}

async function sendMessage(id) {
    var value = document.getElementById("cell-" + id).value
            
    document.getElementById("cell-"+id).style.background = 'white'
            
    console.log("You have selected : ", value)
            
    if(parseInt(value) >= 1 && parseInt(value) <= 9) {    
        row = parseInt(id)/9
        row = Math.floor(row)
        col = parseInt(id)%9
        var userData = {
            "value": parseInt(value),
            "row": parseInt(row),
            "col": parseInt(col)
        }
        conn.send(JSON.stringify(userData))
        conn.onmessage = function (server_data) {
        if(server_data.data == 'violation'){
            document.getElementById("cell-"+id).style.background = 'red'
        } else if(server_data.data == 'win') {
            (async() => {
                const {value: name} = await Swal.fire({
                    title: 'You Won !',
                    imageUrl: 'https://res.cloudinary.com/mayur-cloud/image/upload/c_scale,h_400,w_600/v1597923424/trophy_pyhh2f.jpg',
                    imageWidth: 400,
                    imageHeight: 300,
                    imageAlt: 'Winner',
                    input: 'text',
                    text: 'Enter your Name',
                    inputAttributes: {
                        autocapitalize: 'off'
                    },
                    showCancelButton: false,
                    confirmButtonText: 'Submit',
                    inputValidator: (value) => {
                        if (!value) {
                          return 'Enter your Name'
                        }
                    }
                })
                conn.send(name)
                loadGame(0)
              })();             
        }
                    
    }
    }else {
        document.getElementById("cell-"+id).style.background = 'red'
    }
            
}