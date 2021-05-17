import Quirk, { Redemption } from "./quirk";
import TCP from "./tcp";
import { getString, getInt } from "./env";
import getStatusline from "./statusline";
import getData from "./get-data";
import getCost from "./get-cost";
import validateVimCommand, { vimChangeMode } from "./vim-commands"
import Command, { CommandType } from "./cmd";
import getType from "./get-type";
import validate, { addValidator } from "./validation";
import IrcClient from "./irc";
import { MessageFromPrime, PrimeMessage } from "./irc/prime-commands";
import ProgramWithMe from "./program-with-me";
import bus from "./message-bus";
import { enable } from "./debug";

enable();

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
    const irc = new IrcClient(getString("OAUTH_NAME"), getString("OAUTH_TOKEN"));
    const pwm = new ProgramWithMe();

    addValidator(validateVimCommand);
    addValidator(pwm.validateFunction);

    irc.on("from-theprimeagen", (message: MessageFromPrime) => {
        if (message.type === PrimeMessage.StartYourEngines) {
            pwm.enableProgramWithMe();
        }
        else if (message.type === PrimeMessage.PumpTheBreaks) {
            pwm.disableProgramWithMe();
        }
        else if (message.type === PrimeMessage.PrimeOnly ||
                 message.type === PrimeMessage.FFA) {
            vimChangeMode(message.type);
        }
    });

    quirk.on("message", (data: Redemption): void => {

        const validationResult = validate(data);

        if (!validationResult.success) {

            tcp.write(
                new Command().reset().
                    setStatusLine(validationResult.error).
                    setType(CommandType.StatusUpdate).buffer
            );

            return;
        }

        bus.emit("quirk-message", data);
        tcp.write(
            new Command().reset().
                setCost(getCost(data)).
                setData(getData(data)).
                setStatusLine(getStatusline(data)).
                setType(getType(data)).buffer
        );

    });
}

run();
