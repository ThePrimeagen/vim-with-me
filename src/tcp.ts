import * as net from "net";
import { EventEmitter } from "events";

export default class TCPSocket extends EventEmitter {
    private connections: net.Socket[] = [];

    constructor(port: number) {
        super();

        const server = net.createServer(async (socket: net.Socket) => {
            console.log("Got a connection");

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

    write(data: Buffer | string): Promise<void> {
        return new Promise((res, rej) => {
            const d = TCPSocket.createEnv(data);
            console.log(d);

            this.connections.forEach(c => {
                try {
                    c.write(d, (e: Error | undefined) => {
                        if (e) {
                            this.emit("connection-error", e);
                            rej(e);
                        } else {
                            res();
                        }
                    });
                } catch (e) {
                    this.emit("connection-error", e);
                    rej(e);
                }
            });
        });
    }

    static createEnv(data: Buffer | string): Buffer {
        // using to much data copying... stop
        const d = data instanceof Buffer ? data : Buffer.from(data);
        const len = d.length;
        const env = Buffer.alloc(len + 2);
        env.writeUInt16BE(len, 0);
        d.copy(env, 2);

        return env;
    }
};

