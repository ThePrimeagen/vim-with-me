# Websocket distributor

This is completely dumb and just forwards messages to all downstreams.

Auth is done by the downstreams which will disconnect a socket if it fails to
auth, which will cause this to disconnect the upstream.

## Usage

Pass each `server:port` downstream on the command line.

