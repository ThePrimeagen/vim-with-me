import { CommandType } from "./cmd";
import { Redemption } from "./quirk";

const map: {[key: string]: CommandType} = {
    "asdf": CommandType.SystemCommand,
    "Xrandr": CommandType.SystemCommand,
    "Giveaway FEM": CommandType.GiveawayEnter,
    "VimInsert": CommandType.VimInsert,
    "VimAfter": CommandType.VimAfter,
    "ProgrammWithMeEnter": CommandType.ProgramWithMeEnter,
};

export default function getType(data: Redemption): CommandType {

    // I could of used better string manipulation, but instead, you don't
    // deserve that level of nice programmingh.  Not on a side project.  You
    // must endure the crap
    if (~data.rewardName.indexOf("Vim Command")) {
        return CommandType.VimCommand;
    }

    return map[data.rewardName];
}

