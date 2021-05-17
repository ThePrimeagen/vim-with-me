import getType from "../get-type";
import { CommandType } from "../cmd";
import { ValidationResult } from "../validation";
import { PrimeMessage } from "../irc/prime-commands";
import { Redemption } from "../quirk";

type TieredCommand = {
    name: string | ((input: string) => boolean),
    getValue: (input: string) => string;
}

function getCommand(input: string, commands: TieredCommand[]) {
    let out: string | null = null;
    for (let i = 0; !out && i < commands.length; ++i) {
        const c = commands[i];
        if (typeof c.name === "string" && c.name === input ||
            typeof c.name === "function" && c.name(input)) {
            out = commands[i].getValue(input);
        }
    }

    return out;
}

// T1 commands? -- what is the cost?
// 180 bones
const t1Commands = [
    "h",
    "j",
    "k",
    "l",
]

// T2 commands?
const t2Commands = [
    "o",
    "O",
    "~",
    "gg",
    "G",
    "zt",
    "zz",
    "zh",
    "zl",
    "zb",
    "H",
    "L",
    "0",
    "$",
    "_",
    ":Sex!",
    "<<",
    ">>",
    "V",
    "v",
    "A",
    "I",
];

const t2NamedCommands: TieredCommand[] = [{
    name: "random",
    getValue() {
        return t2Commands[Math.floor(Math.random() * t2Commands.length)];
    }
}, {
    name: "C-a",
    getValue() {
        return "a";
    }
}, {
    name: "C-d",
    getValue() {
        return "d";
    }
}, {
    name: "C-u",
    getValue() {
        return "u";
    }
}, {
    name(input: string) {
        return t1Commands.includes(input[input.length - 1]) &&
            !isNaN(+input.substring(0, input.length - 1));
    },
    getValue(input: string) {
        return input;
    }
}];

// T3 commands?
const t3Commands = [
    "p",
    "P",
    "mA",
    "'A",
    "yiW",
    "viW",
    "dd",
    "dw",
    "diw",
    "diW",
    "cw",
    "ciw",
    "ciW",
    "cc",
    "C",
    "D",
    "S",
];

const t4Commands = [
    "u", // <--- ohhhhhhh
    "cip",
    "cap",
    "ci[",
    "ca[",
    "ci{",
    "ca{",
    "ci(",
    "ca(",
    "dip",
    "dap",
    "di[",
    "da[",
    "di{",
    "da{",
    "di(",
    "da(",
];

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
    const isPWM = name.includes("PWM")

    // BECAUSE I CAN
    if ((isPWM || ~name.indexOf("Tier 1")) && t1Commands.includes(input)) {
        return "";
    }
    if ((isPWM || ~name.indexOf("Tier 2"))) {
        if (t2Commands.includes(input)) {
            return "";
        }
        const cmd = getCommand(input, t2NamedCommands);
        if (cmd) {
            data.userInput = cmd;
            return "";
        }
    }
    if ((isPWM || ~name.indexOf("Tier 3")) && t3Commands.includes(input)) {
        return "";
    }
    if ((isPWM || ~name.indexOf("Tier 4")) && t4Commands.includes(input)) {
        return "";
    }

    return `You cannot use ${data.userInput} at this tier cost.`;
}

// TODO: Think about this better
let mode = PrimeMessage.FFA;
export function vimChangeMode(m: PrimeMessage): void {
    mode = m;
}

export default function validateVimCommand(data: Redemption): ValidationResult {
    if (mode === PrimeMessage.PrimeOnly &&
        data.username !== "ThePrimeagen") {
        return {
            success: false,
            error: `Sorry ${data.username} its prime only mode`,
        };
    }

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


