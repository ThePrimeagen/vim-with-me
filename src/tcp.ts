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
        const d = TCPSocket.createEnv(data);

        this.connections.forEach(c => {
            try {
                c.write(d, (e: Error | undefined) => {
                    if (e) {
                        this.emit("connection-error", e);
                    }
                });
            } catch (e) {
                this.emit("connection-error", e);
            }
        });
    }

    static createEnv(data: Buffer | string): Buffer {
        // using to much data copying... stop
        const d = data instanceof Buffer ? data : Buffer.from(data);
        const len = d.length;
        const env = Buffer.alloc(len + 2);
        env.writeUInt16BE(len, 0);
        env.copy(d, 0, 2);

        return env;
    }
};

