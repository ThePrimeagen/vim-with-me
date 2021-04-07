import * as WebSocket from "ws";

const ws = new WebSocket("ws://vwm.theprimeagen.tv:42069");

ws.on("message", function(data) {
    console.log("Data", data.toString());
});




