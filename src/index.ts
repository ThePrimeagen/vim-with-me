import Quirk, { Redemption } from "./quirk";
import TCP from "./tcp";
import { getString, getInt } from "./env";
import getStatusline from "./statusline";
import getData from "./get-data";
import getCost from "./get-cost";
import validateVimCommand from "./vim-commands"
import primeOnlyCommand, { primeOnlyMode } from "./validation/prime-only"
import Command, { CommandType } from "./cmd";
import getType from "./get-type";
import validate, { addValidator } from "./validation";
import IrcClient from "./irc";
import { MessageFromPrime, PrimeMessage } from "./irc/prime-commands";
import ProgramWithMe from "./program-with-me";
import bus from "./message-bus";
import { enable } from "./debug";
import SystemCommand from "./system-commands";

enable();

const systemCommands: {[key: string]: SystemCommand} = {
    "asdf": new SystemCommand("setxkbmap us", "setxkbmap us real-prog-dvorak", 3000),
    "Xrandr": new SystemCommand("xrandr --output HDMI-0 --brightness 0.05", "xrandr --output HDMI-0 --brightness 1", 5000),
};

async function run(): Promise<void> {

    if (!process.env.QUIRK_TOKEN) {
        console.error("NO ENVIRONMENT (please provide QUIRK_TOKEN)");
        process.exit(1);
    }

    if (!process.env.PORT) {
        console.error("NO ENVIRONMENT (please provide PORT)");
        process.exit(1);
    }

    const quirk = await Quirk.create(getString("QUIRK_TOKEN"));
    const tcp = new TCP(getInt("PORT"));

    // TODO: something....
    // TODO: Maybe we should only emit from this entry point the IRC messages
    // from prime and not from the whole dang irc client...
    // @ts-ignore
    const irc = new IrcClient(getString("OAUTH_NAME"), getString("OAUTH_TOKEN"));
    const pwm = new ProgramWithMe();

    addValidator(primeOnlyCommand);
    addValidator(validateVimCommand);
    addValidator(pwm.validateFunction);

    bus.on("from-theprimeagen", (message: MessageFromPrime) => {
        if (message.type === PrimeMessage.StartYourEngines) {
            pwm.enableProgramWithMe();
        }
        else if (message.type === PrimeMessage.PumpTheBreaks) {
            pwm.disableProgramWithMe();
        }
        else if (message.type === PrimeMessage.PrimeOnly ||
                 message.type === PrimeMessage.FFA) {
            primeOnlyMode(message.type);
        }
    });

    bus.on("system-command", function(command: string) {

        console.log("System Command", command);
        tcp.write(
            new Command().reset().
                setCost(0).
                setData(Buffer.from(command)).
                setStatusLine(``).
                setType(CommandType.SystemCommand).buffer
        );

    });

    quirk.on("message", (data: Redemption): void => {

        const validationResult = validate(data);

        console.log("quirk-message", data, validationResult);
        bus.emit("quirk-message", data, validationResult);

        if (!validationResult.success) {

            tcp.write(
                new Command().reset().
                    setStatusLine(validationResult.error).
                    setType(CommandType.StatusUpdate).buffer
            );

            return;
        }

        const type = getType(data);
        if (type === CommandType.SystemCommand &&
            systemCommands[data.rewardName]) {

            systemCommands[data.rewardName].add();
        } else {
            tcp.write(
                new Command().reset().
                    setCost(getCost(data)).
                    setData(getData(data)).
                    setStatusLine(getStatusline(data)).
                    setType(getType(data)).buffer
            );
        }
    });
}

run();
