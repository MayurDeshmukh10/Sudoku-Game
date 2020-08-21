var input = {
    hours: 0,
    minutes: 0,
    seconds: 0
};

var timestamp = new Date(input.hours, input.minutes, input.seconds);

var interval = 1;

setInterval(function () {
    timestamp = new Date(timestamp.getTime() + interval * 1000);
    input.minutes = timestamp.getMinutes()
    input.hours = timestamp.getHours()
    input.seconds = timestamp.getSeconds()
    document.getElementById('countdown2').innerHTML = timestamp.getHours() + 'h:' + timestamp.getMinutes() + 'm:' + timestamp.getSeconds() + 's';
}, Math.abs(interval) * 1000);