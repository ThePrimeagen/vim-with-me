import Quirk, { Redemption } from "./quirk";
import TCP from "./tcp";
import { getString, getInt } from "./env";
import getStatusline from "./statusline";
import getData from "./get-data";
import getCost from "./get-cost";
import validateVimCommand from "./vim-commands"
import Command, { CommandType } from "./cmd";
import getType from "./get-type";

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

    quirk.on("message", (data: Redemption): void => {

        const validationResult = validateVimCommand(data);

        if (!validationResult.success) {

            tcp.write(
                new Command().reset().
                    setStatusLine(validationResult.error).
                    setType(CommandType.StatusUpdate).buffer
            );

            return;
        }

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
