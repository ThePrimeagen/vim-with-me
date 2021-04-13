import * as net from "net"

const socket = new net.Socket();

socket.on("connect", function() {
    console.log("connected");
});

socket.on("data", function(data: Buffer) {
    console.log("data", data.toString());
});

if (process.env.LOCALHOST) {
    socket.connect(42069);
} else {
    socket.connect(42069, "vwm.theprimeagen.tv");
}
