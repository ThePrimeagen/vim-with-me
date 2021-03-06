import { Redemption } from "./quirk";
import { commandToString, CommandType } from "./cmd";
import getType from "./get-type";

export default function statusLine(data: Redemption, validInput: boolean = true): string {

    const name = data.username;
    if (!validInput) {
        return `Hey, you are *clap emote* fat ${name}`;
    }

    const type = getType(data);
    if (type === CommandType.GiveawayEnter) {
        return `${name}: Thanks for entering the giveaway`;
    }

    else if (type === CommandType.VimInsert) {
        return `${name}: Has inserted ${data.userInput}`;
    }

    else if (type === CommandType.ProgramWithMeEnter) {
        console.log(`statusLine: CommandType ProgramWithMeEnter ${name}`);
        return `${name}: is now PWM`;
    }

    return `${name}: ${commandToString(type).substr(0, 1)} with ${data.userInput}`;
}

