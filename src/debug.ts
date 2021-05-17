import { commandToString, CommandType } from "./cmd";
import getType from "./get-type";
import bus from "./message-bus"
import { Redemption } from "./quirk";

let enabled = false;

bus.on("quirk-message", function(data: Redemption) {
    if (enabled) {
        return;
    }

    const type = getType(data);
    if (type === CommandType.VimAfter || type === CommandType.VimInsert ||
            type === CommandType.VimCommand) {

        console.log(`${data.username}: ${commandToString(getType(data))}: ${data.userInput}`);
    }
});


export function enable() {
    enabled = true;
}

export function disable() {
    enabled = false;
}
