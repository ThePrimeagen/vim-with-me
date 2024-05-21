import { parseFrame } from "./net/frame.js"
import { WS } from "./ws/index.js"

// TODO: provide a url?

/**
 * @param {HTMLElement} el
 */
function run(el) {
    const ws = new WS("ws://localhost:8080/ws")

    ws.onMessage(async function(blob) {
        const buf = new Uint8Array(await blob.arrayBuffer())
        const frame = parseFrame(buf)
        console.log(frame.cmd, frame.data.byteLength)
    })
}

run(document.body)
