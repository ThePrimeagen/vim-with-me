import * as net from "net"

const socket = new net.Socket();

socket.on("connect", function() {
    console.log("connected");
});

socket.on("data", function(data: Buffer) {
    console.log("data", data.toString());
});

socket.connect(42069);
