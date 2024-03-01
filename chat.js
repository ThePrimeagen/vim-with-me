#!/usr/bin/env node

const log = require("./logger");
const tmi = require('tmi.js');

const client = new tmi.Client({
	channels: [ 'theprimeagen' ]
});

client.connect();

client.on('message', (channel, tags, message, self) => {
    log(`message:${tags['display-name']}:${message}`);
});

client.on("cheer", (channel, userstate, message) => {
    log(`bits:${userstate['display-name']}:${userstate.bits}:${message}`);
});

