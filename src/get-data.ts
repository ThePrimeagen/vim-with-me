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
            out = Buffer.from(data.userInput);
            break;
    }

    return out;
}

