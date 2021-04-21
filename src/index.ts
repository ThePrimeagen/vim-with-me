import Quirk, { Redemption } from "./quirk";
import TCP from "./tcp";
import { getString, getInt } from "./env";
import getStatusline from "./statusline";
import Command, { CommandType } from "./cmd";

async function run() {

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

    quirk.on("message", (data: Redemption) => {
        console.log("quirk redemption", data);
        const statusline = getStatusline(data);

        console.log("Status line:", statusline);
        tcp.write(
            new Command().reset().
                setStatusLine(statusline).
                setType(CommandType.StatusUpdate).buffer
        );
    });
}

run();

