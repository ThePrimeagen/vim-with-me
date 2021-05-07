import { Redemption } from "../quirk";
import getType from "../get-type";
import { CommandType } from "../cmd";
import { ValidationResult } from "../validation";

// T1 commands?
const t1Commands = [
    "h",
    "j",
    "k",
    "l",
]

// T2 commands?
const t2Commands = [
    "S",
    "dd",
    "dw",
    "db",
    "w",
    "b",
]

function printCharacters(data: Redemption): void {
    console.log(`${data.username}:`, data.userInput.split("").map(x => x.charCodeAt(0)));
}

function hasBadCharacters(str: string): boolean {
    return Boolean(str.split("").
        map(x => x.charCodeAt(0)).
        filter(x => x < 32 || x > 127).length);
}

function convertEscapeCharacters(str: string): string {
    return str.split("\\n").join("\r").split("\\r").join("\r");
}

function insert(data: Redemption): string {
    printCharacters(data);
    if (hasBadCharacters(data.userInput)) {
        return "You cannot use 32 < ascii > 127";
    }

    data.userInput = convertEscapeCharacters(data.userInput);

    if (data.userInput.length > 5) {
        data.userInput = data.userInput.substr(0, 5);
    }

    return "";
};

function vimCommand(data: Redemption): string {

    // starting simple
    const input = data.userInput;
    const name = data.rewardName;

    // BECAUSE I CAN
    if (~name.indexOf("Tier 1") && t1Commands.includes(input)) {
        return "";
    } else if (~name.indexOf("Tier 2") && t2Commands.includes(input)) {
        return "";
    } else if (~name.indexOf("Tier 3")) {
        return "";
    }

    return `You cannot use ${data.userInput} at this tier cost.`;
}

export default function validateVimCommand(data: Redemption): ValidationResult {
    const type = getType(data);

    let error: string = "";
    switch (type) {
    case CommandType.VimCommand:
        error = vimCommand(data);
        break;
    case CommandType.VimInsert:
    case CommandType.VimAfter:
        error = insert(data);
        break;
    }

    if (error === "") {
        return {
            success: true
        };
    }

    return {
        success: false,
        error,
    };
}


