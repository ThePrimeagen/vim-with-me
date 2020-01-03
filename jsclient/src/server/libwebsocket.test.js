const WebSockets = require('ws');

const server = new WebSockets.Server({
    path: "/",
    port: 8080
});

// @ts-ignore
server.on('connection', ws => {
    let intervalId = 0;
    let count = 0;

    console.log("Connection!!!");

    ws.on('message', d => {
        console.log("Message?", d);
        ws.send(`hello world ${++count}`);
    });

    ws.on('close', () => {
        clearInterval(intervalId);
    });
    ws.send(`hello world ${++count}`);
});


