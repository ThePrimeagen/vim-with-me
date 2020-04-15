import WebSocket from 'ws';

import * as cp from 'child_process';

function wait(ms: number) {
    return new Promise(res => {
        setTimeout(res, ms);
    });
}

// Commands.  E Z P Z
// Attach to neovim process
async function run() {
    require('neovim/scripts/nvim').then((nvim) => {
        nvim.command('vsp');
    });
}

if (require.main === module) {
    run();
}


