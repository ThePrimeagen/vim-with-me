import Quirk, { Redemption } from "./quirk";
import TCP from "./tcp";
import { getString, getInt } from "./env";
import getStatusline from "./statusline";
import getType from "./get-type";
import getData from "./get-data";
import getCost from "./get-cost";
import Command, { CommandType, validInput } from "./cmd";

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
        console.log("quirk redemption", data);
        const type = getType(data);

        //if (type === CommandType.VimCommand && !validInput(data.userInput)) {
        if (type === CommandType.VimCommand && !validInput(data.userInput)) {
            tcp.write(
                new Command().reset().
                    setStatusLine(getStatusline(data, false)).
                    setType(CommandType.StatusUpdate).buffer
            );

            return;
        }

        tcp.write(
            new Command().reset().
                setCost(getCost(data)).
                setData(getData(data)).
                setStatusLine(getStatusline(data)).
                setType(type).buffer
        );

    });
}

run();
