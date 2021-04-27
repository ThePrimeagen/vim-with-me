declare global {
    export type CommandProcessor = (cmd: Command, quirk: Redemption) => string;
}

