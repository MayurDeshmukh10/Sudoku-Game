var conn
function loadGame() {
    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.close = function () {
    console.log("Connection with Server closed")
    }
    conn.onmessage = function (server_data) {
        var message = server_data.data
        console.log("Message len :", message.length)
        var tr;
        var tbl = document.getElementById("grid")
        for (var i = 0; i < message.length; i++) {
            if (i % 9 == 0) {
                console.log("In tr")
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
            // console.log(message[i])

        }
        console.log("Server : ", message)
    }

}
        
window.onload = function () {
    if (window["WebSocket"]) {
        loadGame()
    } else {
        console.log("Your Browser does not websocket")
        alert("Your Browser does not websocket")
    }
}

        
function resetGame() {
    location.reload(); 
}
        
function sendMessage(id) {

    console.log("Time : ",input.minutes)
            
    var value = document.getElementById("cell-" + id).value
            
    document.getElementById("cell-"+id).style.background = 'white'
            
    console.log("You have selected : ", value)
            
    if(parseInt(value) >= 1 && parseInt(value) <= 9) {    
        row = parseInt(id)/9
        row = Math.floor(row)
        col = parseInt(id)%9
        conn.send(value + ',' + row + ',' + col)
        conn.onmessage = function (server_data) {
        console.log(server_data.data)
        if(server_data.data == 'violation'){
            document.getElementById("cell-"+id).style.background = 'red'
        } else if(server_data.data == 'win') {
            Swal.fire({
            title: 'You Won !',
            imageUrl: 'https://res.cloudinary.com/mayur-cloud/image/upload/c_scale,h_400,w_600/v1597923424/trophy_pyhh2f.jpg',
            imageWidth: 400,
            imageHeight: 300,
            imageAlt: 'Winner',
            })
        }
                    
    }
    }else {
        document.getElementById("cell-"+id).style.background = 'red'
    }
            
}