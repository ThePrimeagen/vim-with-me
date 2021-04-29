import { CommandType } from "./cmd";
import { Redemption } from "./quirk";

const map: {[key: string]: CommandType} = {
    "Vim Command": CommandType.VimCommand,
    "ASDF": CommandType.ASDF,
    "Xrandr": CommandType.Xrandr,
    "Giveaway FEM": CommandType.GiveawayEnter,
};

export default function getType(data: Redemption): CommandType {
    if (~data.rewardName.indexOf("Vim Command")) {
        return CommandType.VimCommand;
    }
    return map[data.rewardName];
}

