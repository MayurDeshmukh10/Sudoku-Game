var input = {
    hours: 0,
    minutes: 0,
    seconds: 0
};
var stopwatch

function timer() {
    var timestamp = new Date(input.hours, input.minutes, input.seconds);
    
    var interval = 1;
    
    stopwatch = setInterval(function () {
        timestamp = new Date(timestamp.getTime() + interval * 1000);
        input.minutes = timestamp.getMinutes()
        input.hours = timestamp.getHours()
        input.seconds = timestamp.getSeconds()
        document.getElementById('countdown2').innerHTML = input.hours + 'h:' + input.minutes + 'm:' + input.seconds + 's';
    }, Math.abs(interval) * 1000);

}

function resetTimer(){
    input.seconds = 0
    input.minutes = 0
    input.hours = 0
    clearInterval(stopwatch)
    document.getElementById('countdown2').innerHTML = '0h:0m:0s';
}

