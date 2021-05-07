import { EventEmitter } from "events";
import { IrcTags } from ".";

export enum PrimeMessage {
    StartYourEngines = 1,
    PumpTheBreaks = 2,
}
export type MessageFromPrime = {
    type: PrimeMessage,
}

export default function primeCommands(emitter: EventEmitter, tags: IrcTags, message: string): void {
    if (tags["display-name"] === "ThePrimeagen") {
        // TODO: He wont remember this ever
        // PLEASE READ THIS ON STREAM YOU 5HEAD
        if (message === "!start-program-with-me") {
            emitter.emit("from-theprimeagen", {
                type: PrimeMessage.StartYourEngines
            });
        } else if (message === "!stop-program-with-me") {
            emitter.emit("from-theprimeagen", {
                type: PrimeMessage.PumpTheBreaks
            });
        }
    }
}

