#!/usr/bin/env node

const tmi = require('tmi.js');

const client = new tmi.Client({
	channels: [ 'theprimeagen' ]
});

client.connect();

client.on('message', (channel, tags, message, self) => {
	// "Alca: Hello, World!"
	console.log(`${tags['display-name']}:${message}`);
});



