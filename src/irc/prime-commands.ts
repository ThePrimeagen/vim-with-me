import { EventEmitter } from "events";
import { IrcTags } from ".";

export enum PrimeMessage {
    StartYourEngines = 1,
    PumpTheBreaks = 2,
    PrimeOnly = 3,
    FFA = 4,
}

const toStringMap = new Map<PrimeMessage, string>([
    [PrimeMessage.StartYourEngines, "StartYourEngines"],
    [PrimeMessage.PumpTheBreaks, "PumpTheBreaks"],
    [PrimeMessage.PrimeOnly, "PrimeOnly"],
    [PrimeMessage.FFA, "FFA"],
]);

export type MessageFromPrime = {
    type: PrimeMessage,
}

export function toStringMessageFromPrime(message: MessageFromPrime): string {
    return toStringMap.get(message.type);
}

const msgToEmit: {[key: string]: PrimeMessage} = {
    "!start-program-with-me": PrimeMessage.StartYourEngines,
    "!stop-program-with-me": PrimeMessage.PumpTheBreaks,
    "!prime-on": PrimeMessage.PrimeOnly,
    "!prime-off": PrimeMessage.FFA,
}

export default function primeCommands(emitter: EventEmitter, tags: IrcTags, message: string): void {
    if (tags["display-name"] === "ThePrimeagen") {

        const type = msgToEmit[message];
        if (type) {
            emitter.emit("from-theprimeagen", {
                type,
            });
        }
    }
}

