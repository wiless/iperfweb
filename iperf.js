console.log('Welcome to iperf script...loaded..');
var ws;
var statuspos = false;
function startstop() {
    console.log("AM I supposed to toggle");
    if (statuspos == false) {
        statuspos = true;
        console.log("Writing something to socket");
        ws.send(JSON.stringify({ Status: "start" }));
    }else {
        statuspos = false;
        console.log("Asking to stop");
        ws.send(JSON.stringify({ Status: "" }));   
    }
}
function onClosefn() {
    console.log("Socket is closed");
}
function onloadfn() {
    console.log("Socket is open");

}
function onmessagefn(jdata) {
    console.log("Throughput " + jdata.Speed + " , " + jdata.Unit);
    // setInterval(function () {
    data.setValue(0, 1, Math.round(jdata.Speed));
    chart.draw(data, optionsLocal);
    // }, 1000);

}
function initialize() {
    drawChart();
    ws = openWS();
}
function openWS() {
    wsurl = "ws://" + window.location.hostname + ":" + window.location.port + "/ws";
    var socket = new WebSocket(wsurl);
    socket.onopen = function (e) {
        console.log("<p>Socket is open</p>");
        onloadfn();
    };
    socket.onmessage = function (e) {
        //console.log(e);
        jdata = JSON.parse(e.data);
        onmessagefn(jdata);
    };
    socket.onclose = function () {
        console.log("Socket closed");
        // container.append("<p>Socket closed</p>");
        onclose();
    };
    return socket;
}

var data, optionsLAN,optionsLocal;
var chart;
function drawChart() {

    data = google.visualization.arrayToDataTable([
        ['Label', 'Value'],
        ['DL Mbps', 0]
    ]);


    //   data = google.visualization.arrayToDataTable([
    //     ['Label', 'Value'],
    //     ['DL Mbps', 0],
    //     ['CPU', 55],
    //     ['Network', 68]
    // ]);
    optionsLocal = {
        width: 400, height: 400,
        redFrom: 10000, redTo: 15000,
        yellowFrom: 5000, yellowTo: 10000,
        minorTicks: 500,
        min: 100, max: 15000
    };

    optionsLAN = {
        width: 400, height: 400,
        redFrom: 100, redTo: 200,
        yellowFrom: 50, yellowTo: 100,
        minorTicks: 10,
        min: 0, max: 200
    };

    chart = new google.visualization.Gauge(document.getElementById('chart_div'));

    chart.draw(data, optionsLocal);

    // setInterval(function () {
    //     data.setValue(0, 1, 40 + Math.round(60 * Math.random()));
    //     chart.draw(data, options);
    // }, 1000);

    // setInterval(function () {
    //     data.setValue(1, 1, 40 + Math.round(60 * Math.random()));
    //     chart.draw(data, options);
    // }, 1500);

    // setInterval(function () {
    //     data.setValue(2, 1, 60 + Math.round(20 * Math.random()));
    //     chart.draw(data, options);
    // }, 500);
}
//  function createGauge(elem,mn,mx) {
//     var gaugePS = new RadialGauge({
//         renderTo: elem,
//         width: 270,
//         height: 270,
//         units: 'Mbps',
//         title: "DL",
//         value:0,
//         minValue: mn,
//         maxValue: mx,
//         majorTicks: [
//             25,50,75,100
//         ],
//         minorTicks: 2,
//         ticksAngle: 270,
//         startAngle: 45, 
//         strokeTicks: true,
//         highlights: [
//             { from: 80, to: 100, color: 'rgba(225, 7, 23, 0.75)' }
//         ],
//         valueInt: 1,
//         valueDec: 1,
//         colorPlate: "#fff",
//         colorMajorTicks: "#686868",
//         colorMinorTicks: "#686868",
//         colorTitle: "#000",
//         colorUnits: "#000",
//         colorNumbers: "#686868",
//         valueBox: true,
//         colorValueText: "#000",
//         colorValueBoxRect: "#fff",
//         colorValueBoxRectEnd: "#fff",
//         colorValueBoxBackground: "#fff",
//         colorValueBoxShadow: false,
//         colorValueTextShadow: false,
//         colorNeedleShadowUp: true,
//         colorNeedleShadowDown: false,
//         colorNeedle: "rgba(200, 50, 50, .75)",
//         colorNeedleEnd: "rgba(200, 50, 50, .75)",
//         colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
//         colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
//         borderShadowWidth: 0,
//         borders: true,
//         borderInnerWidth: 0,
//         borderMiddleWidth: 0,
//         borderOuterWidth: 5,
//         colorBorderOuter: "#fafafa",
//         colorBorderOuterEnd: "#cdcdcd",
//         needleType: "arrow",
//         needleWidth: 2,
//         needleCircleSize: 7,
//         needleCircleOuter: true,
//         needleCircleInner: false,
//         animationDuration: 500,
//         animation:true,
//         animationRule: "dequint",
//         fontNumbers: "Verdana",
//         fontTitle: "Verdana",
//         fontUnits: "Verdana",
//         fontValue: "Led",
//         fontValueStyle: 'italic',
//         fontNumbersSize: 20,
//         fontNumbersStyle: 'italic',
//         fontNumbersWeight: 'bold',
//         fontTitleSize: 24,
//         fontUnitsSize: 22,
//         fontValueSize: 50,
//         animatedValue: true
//     });
//  }
