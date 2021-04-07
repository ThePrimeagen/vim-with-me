import * as ws from "ws";

const server = new ws.Server({
    port: +process.env.PORT
});

server.on('connection', function connection(ws) {
    ws.send('Hello, world1');
    ws.close();
});

