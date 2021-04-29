import { Redemption } from "./quirk";
import { CommandType } from "./cmd";
import getType from "./get-type";

export default function getData(data: Redemption): null | Buffer {
    const type = getType(data);
    let out: Buffer | null = null;

    // TODO: We did this because we are engineers and this is clearly the most
    // bestest way to be abstracted from the complications as types grow
    // instead of using a map.  AMIRIGHT??? SWITCH STATEMENTS ARE BAE
    switch (type) {
        case CommandType.VimCommand:
            // TODO: Probably need to think about this one... huh?
            out = Buffer.from(`norm ${data.userInput}`);
            break;
        case CommandType.GiveawayEnter:
            out = Buffer.from(data.username);
            break;
    }

    return out;
}

