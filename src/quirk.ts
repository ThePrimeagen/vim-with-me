import { EventEmitter } from "events";

import fetch from "node-fetch";
import * as WebSocket from "ws";

type TwitchRedemption = {
    source: "TWITCH_EVENTSUB" | "TWITCH_PUBSUB",
    type: "TWITCH_CHANNEL_REWARD"
    data: {
        user_name: string,
        user_input: string,
        reward: {
            title: string,
            cost: number,
        }
    },
}

export type Redemption = {
    username: string,
    userInput: string,
    rewardName: string,
    cost: number,
}

async function connectToQuirk(token: string): Promise<WebSocket> {
    const domain = "websocket.quirk.gg";
    const body = {
        access_token: token,
    };

    const response = await fetch(`https://${domain}/token`, {
        body: JSON.stringify(body),
        headers: { "Content-Type": "application/json" },
        method: "POST",
    });

    const data = await response.json();

    const { access_token } = data;

    return new WebSocket(
        `wss://websocket.quirk.gg?access_token=${access_token}`
    );
}

export default class Quirk extends EventEmitter {
    static create(token: string): Promise<Quirk> {
        return new Promise(async (res, rej) => {
            const socket = await connectToQuirk(token);

            function open() {
                res(new Quirk(socket));
                socket.off("open", open);
                socket.off("error", error);
            };

            function error(e: Error) {
                rej(e);
                socket.off("open", open);
                socket.off("error", error);
            };

            socket.on("error", error);
            socket.on("open", open);
        });
    }

    private constructor(socket: WebSocket) {
        super();

        socket.on("message", (data: Buffer) => {
            const parsedData: TwitchRedemption = JSON.parse(data.toString());
            try {
                if (parsedData.source === "TWITCH_EVENTSUB" &&
                    parsedData.type === "TWITCH_CHANNEL_REWARD" &&
                    parsedData.data.reward && parsedData.data.reward.title) {
                    this.emit("message", {
                        username: parsedData.data.user_name,
                        userInput: parsedData.data.user_input,
                        rewardName: parsedData.data.reward.title,
                        cost: +parsedData.data.reward.cost,
                    });
                }
            } catch (e) {
                console.log("ERRROR", e.message);
            }
        });

        socket.on("error", this.emit.bind(this, "error"));
    }
}

async function test() {
    const q = await Quirk.create(String(process.env.QUIRK));
    q.on("message", (e) => console.log(JSON.stringify(e, null, 4)));
}

if (require.main === module) {
    test();
}
