import IrcClient from "./irc";
import { getString } from "./env";

/**
const socket = new net.Socket();

socket.on("connect", function() {
    console.log("connected");
});

socket.on("data", function(data: Buffer) {
    console.log("data", data.toString());
});

if (process.env.LOCALHOST) {
    socket.connect(42069);
} else {
    socket.connect(42069, "vwm.theprimeagen.tv");
}


*/

const client = new IrcClient(getString("OAUTH_NAME"), getString("OAUTH_TOKEN"));

client.on("connect", function() {
    console.log("On Connected baybee");
});

client.on("from-theprimeagen", function(m) {
    console.log("MESSAGE", m);
});
