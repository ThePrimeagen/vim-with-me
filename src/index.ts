import * as net from "net";

import Quirk, { Redemption } from "./quirk";
import { getString } from "./env";

async function run() {
    if (!process.env.QUIRK_TOKEN) {
        console.error("NO ENVIRONMENT (please provide QUIRK_TOKEN)");
        process.exit(1);
    }

    const quirk = await Quirk.create(getString("QUIRK_TOKEN"));
    const connections: net.Socket[] = [];
    const server = net.createServer(function(socket): void {
        connections.push(socket);
        console.log("New Connection!!!");

        socket.on("close", function(): void {
            connections.splice(connections.indexOf(socket), 1);
        });
    });

    quirk.on("message", (data: Redemption) => {
        console.log("quirk redemption", data);
        connections.forEach(c => {
            c.write("Redemption: " + JSON.stringify(data));
        });
    });


    server.on("error", function(e) {
        console.error("Server Error: ", e);
        process.exit(1);
    });

    console.log("Listening too", process.env.PORT);
    server.listen(process.env.PORT);
}

run();

