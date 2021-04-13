import * as net from "net";
import { EventEmitter } from "events";

export default class TCPSocket extends EventEmitter {
    private connections: net.Socket[] = [];

    constructor(port: number) {
        super();

        const server = net.createServer((socket: net.Socket) => {
            this.connections.push(socket);
            socket.on("close", () => {
                this.connections.splice(this.connections.indexOf(socket), 1);
            });
        });

        server.on("error", this.emit.bind(this, "error"));
        server.listen(port, () => {
            this.emit("listening");
        });
    }

    write(data: Buffer | string): void {
        this.connections.forEach(c => {
            try {
                c.write(data, (e: Error | undefined) => {
                    if (e) {
                        this.emit("connection-error", e);
                    }
                });
            } catch (e) {
                this.emit("connection-error", e);
            }
        });
    }
};

