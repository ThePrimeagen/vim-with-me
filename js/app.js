import { asciiPixel, createDecodeFrame, expand, pushFrame } from "./decode/frame.js"
import { parseFrame } from "./net/frame.js"
import { WS } from "./ws/index.js"

// TODO: provide a url?

/**
 * @param {HTMLElement} el
 */
function run(el) {
    const ws = new WS("ws://localhost:8080/ws")

    const decodeFrame = createDecodeFrame()
    ws.onMessage(async function(buf) {
        const frame = parseFrame(buf)

        pushFrame(decodeFrame, frame)
        expand(decodeFrame)

        const data = asciiPixel(decodeFrame)
        console.log(data.slice(0, 5))
    })
}

run(document.body)
